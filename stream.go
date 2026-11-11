package status302

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"maps"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// StreamAuth selects how a Stream authenticates its WebSocket upgrade.
type StreamAuth string

const (
	// StreamAuthBearer sends the API key as `Authorization: Bearer` on the upgrade request.
	StreamAuthBearer StreamAuth = "bearer"
	// StreamAuthTicket connects with `?ticket=` from StreamOptions.GetTicket,
	// fetched before every (re)connect.
	StreamAuthTicket StreamAuth = "ticket"
)

// StreamEncoding is the wire encoding of server → client messages.
type StreamEncoding string

const (
	// StreamEncodingJSON: every server message is a JSON text frame (the default).
	StreamEncodingJSON StreamEncoding = "json"
	// StreamEncodingMsgpack: every server message is one binary MessagePack
	// frame (`?encoding=msgpack`), decoded into the same messages as JSON.
	// Client messages stay JSON text.
	StreamEncodingMsgpack StreamEncoding = "msgpack"
)

// StreamState is the connection state of a Stream (see StreamOptions.OnStateChange).
type StreamState string

const (
	// StreamIdle: nothing subscribed yet, so not connected.
	StreamIdle StreamState = "idle"
	// StreamConnecting: the first connection is being opened.
	StreamConnecting StreamState = "connecting"
	// StreamOpen: connected; the subscriptions are (re)sent.
	StreamOpen StreamState = "open"
	// StreamReconnecting: the connection failed or ended; waiting out the backoff or dialing again.
	StreamReconnecting StreamState = "reconnecting"
	// StreamClosed: ended for good (Close, a final close code or a refused ticket).
	StreamClosed StreamState = "closed"
)

// Close codes the server uses (https://status302.com/docs/stream/#close-codes).
const (
	// CloseGoingAway: idle timeout or a failed write on the server. Reconnects.
	CloseGoingAway = 1001
	// CloseServiceRestart: the server is restarting. Reconnects and resumes with `since`.
	CloseServiceRestart = 1012
	// CloseSlowConsumer: the connection's send buffer overflowed. Reconnects and resumes with `since`.
	CloseSlowConsumer = 4008
	// CloseUnauthorized: credentials missing, invalid or revoked. Final.
	CloseUnauthorized = 4401
	// CloseForbidden: the plan does not allow streaming, or the API key lacks the `stream` scope. Final.
	CloseForbidden = 4403
	// CloseReconnect: role, org, plan or key scopes changed. Reconnects at once
	// with fresh credentials and resumes.
	CloseReconnect = 4409
	// CloseTooManyConnections: the org's stream connection quota is used up.
	// Reported through StreamOptions.OnError, then retried with backoff.
	CloseTooManyConnections = 4429
)

const (
	// closeIdle is the code the client closes with when nothing arrived for IdleTimeout.
	closeIdle = 4000
	// subscribeRetryDelay is the wait before resending a subscribe refused as
	// `unavailable` or `rate_limited`.
	subscribeRetryDelay = time.Second
	// maxMsgpackDepth bounds the nesting of a MessagePack frame (as encoding/json does).
	maxMsgpackDepth = 10000
)

// coalescedTopic matches the topics whose messages are the latest state only:
// `since` does not apply, a resubscribe after a reconnect reloads the snapshot.
var coalescedTopic = regexp.MustCompile(`^(event\.\d+\.(kpis|map)|catalog\.rows|catalog\.movers)$`)

var jsonNull = json.RawMessage("null")

// ErrStreamClosed is returned by Stream.Next once the stream has ended (after
// Close, or after the final *StreamClosedError was returned), and by
// Subscription.Ready when the stream ends before the server confirmed.
var ErrStreamClosed = errors.New("status302: stream is closed")

// StreamClosedError ends a stream for good: the server closed it with 4401 or
// 4403, or a ticket was refused (Code 4401, Reason "ticket refused"; the
// *APIError is wrapped). Stream.Next returns it once. StreamOptions.OnError
// also receives it for 4429, which is retried.
type StreamClosedError struct {
	Code   int
	Reason string
	err    error
}

func (e *StreamClosedError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("status302: stream closed (%d)", e.Code)
	}
	return fmt.Sprintf("status302: stream closed (%d %s)", e.Code, e.Reason)
}

// Unwrap returns the cause, e.g. the *APIError of a refused ticket.
func (e *StreamClosedError) Unwrap() error { return e.err }

// SubscriptionError is a subscribe request the server refused
// (`forbidden_topic`, `unknown_topic`, `bad_request`, `subscription_limit`).
type SubscriptionError struct {
	Topic   string
	Code    string
	Message string
}

func (e *SubscriptionError) Error() string {
	return fmt.Sprintf("status302: %s: %s: %s", e.Topic, e.Code, e.Message)
}

// Message is a message delivered by Stream.Next: a *DataMessage, a
// *ResyncMessage or an *ErrorMessage.
type Message interface {
	// StreamTopic is the topic the message belongs to.
	StreamTopic() string
}

// DataMessage is a topic message. Seq increases per topic; Data is the topic
// payload (re-encoded as JSON with the msgpack encoding).
type DataMessage struct {
	Topic string
	Seq   int64
	TS    string
	Data  json.RawMessage
}

// StreamTopic is the message's topic.
func (m *DataMessage) StreamTopic() string { return m.Topic }

// Decode unmarshals Data into v.
func (m *DataMessage) Decode(v any) error { return json.Unmarshal(m.Data, v) }

// ResyncMessage asks to reload Snapshot (a REST path, see Stream.FetchSnapshot).
// Reason is `gap` (`since` is older than the replay window, or the server's
// seq counter restarted: seqs start over after it), `backpressure` (messages
// were dropped) or `reconnect` (a coalesced topic was resubscribed after a
// reconnect).
type ResyncMessage struct {
	Topic    string
	Reason   string
	Snapshot string
}

// StreamTopic is the message's topic.
func (m *ResyncMessage) StreamTopic() string { return m.Topic }

// ErrorMessage reports that the server refused the subscription to Topic; it
// is dropped and its Subscription.Ready returns a *SubscriptionError.
type ErrorMessage struct {
	Topic   string
	Code    string
	Message string
}

// StreamTopic is the message's topic.
func (m *ErrorMessage) StreamTopic() string { return m.Topic }

// StreamOptions configures Client.Stream. The zero value streams JSON from
// the client's base URL with the client's API key.
type StreamOptions struct {
	// APIKey is sent with bearer auth (default: the client's key).
	APIKey string
	// Auth is StreamAuthBearer (server-side) or StreamAuthTicket. Default:
	// bearer when an API key is set, else ticket.
	Auth StreamAuth
	// GetTicket returns a fresh single-use ticket; called before every
	// (re)connect with ticket auth. Default: the client's CreateStreamTicket.
	GetTicket func(ctx context.Context) (string, error)
	// URL is the WebSocket URL (default StreamURL(client.BaseURL())).
	URL string
	// Encoding of server messages (default StreamEncodingJSON).
	Encoding StreamEncoding
	// ReconnectDelay is the first reconnect delay (default 500 ms); it
	// doubles per failed attempt, with full jitter.
	ReconnectDelay time.Duration
	// MaxReconnectDelay caps the reconnect delay (default 30 s).
	MaxReconnectDelay time.Duration
	// IdleTimeout reconnects when nothing arrives for this long (default 60 s;
	// the server sends a heartbeat every 20 s). It also bounds the handshake.
	IdleTimeout time.Duration
	// HTTPClient dials the WebSocket (default http.DefaultClient).
	HTTPClient *http.Client
	// OnStateChange receives every state change, in order, from a goroutine
	// of the stream.
	OnStateChange func(StreamState)
	// OnError receives non-final connection problems, in order, from a
	// goroutine of the stream: 4429 closes (*StreamClosedError), ticket
	// errors, malformed messages and server errors that answer no pending
	// request. Final errors are returned by Stream.Next instead.
	OnError func(error)
}

// SubscribeOptions configures Stream.Subscribe.
type SubscribeOptions struct {
	// Filter narrows the topic, e.g. {"country": ["US"], "minPrice": 10000} for `sales`.
	Filter map[string]any
	// Since resumes after this seq (replayable topics). Later reconnects
	// resume from the last delivered seq automatically.
	Since *int64
}

// Stream is the realtime stream client (protocol: https://status302.com/docs/stream/reference/). One
// WebSocket carries every subscription; the stream connects on the first
// Subscribe, reconnects with backoff, resumes every replayable topic with
// `since` from the last delivered seq, drops duplicates by seq, and delivers
// a *ResyncMessage when the server cannot replay a gap (or a coalesced topic
// was resubscribed after a reconnect):
//
//	stream := client.Stream(status302.StreamOptions{})
//	defer stream.Close()
//	stream.Subscribe("sales", status302.SubscribeOptions{Filter: map[string]any{"country": []string{"US"}}})
//	for msg, err := range stream.Messages(ctx) {
//		if err != nil {
//			return err // a *StreamClosedError (4401, 4403, refused ticket) or ctx.Err()
//		}
//		switch m := msg.(type) {
//		case *status302.DataMessage:
//			fmt.Println(m.Topic, m.Seq, string(m.Data))
//		case *status302.ResyncMessage:
//			snapshot, err := stream.FetchSnapshot(ctx, m)
//			…
//		case *status302.ErrorMessage:
//			log.Printf("%s refused: %s", m.Topic, m.Code)
//		}
//	}
//
// Messages queue up without bound until Next takes them. A Stream is safe
// for concurrent use.
type Stream struct {
	client      *Client
	url         string
	header      http.Header                               // bearer credentials
	getTicket   func(ctx context.Context) (string, error) // nil with bearer auth
	encoding    StreamEncoding
	baseDelay   time.Duration
	maxDelay    time.Duration
	idleTimeout time.Duration
	httpClient  *http.Client
	onState     func(StreamState)
	onError     func(error)
	configErr   error // returned by Next instead of connecting

	ctx       context.Context // ends with the stream: bounds tickets, dials and reads
	cancel    context.CancelFunc
	done      chan struct{} // closed when the stream ends
	wake      chan struct{} // a message was queued
	eventWake chan struct{} // a callback event was queued

	mu            sync.Mutex
	state         StreamState
	started       bool
	closed        bool // ended: Close or a final error
	userClosed    bool
	final         error // returned once by Next after the queued messages
	subs          map[string]*Subscription
	pending       map[string]*Subscription // request id → subscription, on the current connection
	conn          *streamConn              // nil while not connected
	connectedOnce bool
	nextID        int
	queue         []Message // queue[head:] awaits Next
	head          int
	events        []streamEvent
}

// streamConn is one WebSocket connection and its outgoing messages.
type streamConn struct {
	ws   *websocket.Conn
	out  [][]byte // guarded by Stream.mu
	wake chan struct{}
}

// streamEvent is a queued OnStateChange (err == nil) or OnError call.
type streamEvent struct {
	state StreamState
	err   error
}

// clientMessage is a subscribe or unsubscribe request.
type clientMessage struct {
	Op     string          `json:"op"`
	ID     string          `json:"id"`
	Topic  string          `json:"topic"`
	Since  *int64          `json:"since,omitempty"`
	Filter json.RawMessage `json:"filter,omitempty"`
}

// serverMessage is the union of the server's messages.
type serverMessage struct {
	Type     string          `json:"type"`
	ID       string          `json:"id"`
	Topic    string          `json:"topic"`
	Seq      int64           `json:"seq"`
	TS       string          `json:"ts"`
	Data     json.RawMessage `json:"data"`
	Reason   string          `json:"reason"`
	Snapshot string          `json:"snapshot"`
	Code     string          `json:"code"`
	Message  string          `json:"message"`
}

// Stream returns a realtime stream client; it connects on the first
// Subscribe. Close it when done.
func (c *Client) Stream(opts StreamOptions) *Stream {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Stream{
		client:      c,
		encoding:    cmp.Or(opts.Encoding, StreamEncodingJSON),
		baseDelay:   cmp.Or(opts.ReconnectDelay, 500*time.Millisecond),
		maxDelay:    cmp.Or(opts.MaxReconnectDelay, 30*time.Second),
		idleTimeout: cmp.Or(opts.IdleTimeout, 60*time.Second),
		httpClient:  opts.HTTPClient,
		onState:     opts.OnStateChange,
		onError:     opts.OnError,
		ctx:         ctx,
		cancel:      cancel,
		done:        make(chan struct{}),
		wake:        make(chan struct{}, 1),
		eventWake:   make(chan struct{}, 1),
		state:       StreamIdle,
		subs:        map[string]*Subscription{},
		pending:     map[string]*Subscription{},
	}
	var errs []error
	s.url = cmp.Or(opts.URL, StreamURL(c.baseURL))
	if u, err := url.Parse(s.url); err != nil || !slices.Contains([]string{"ws", "wss", "http", "https"}, u.Scheme) {
		errs = append(errs, fmt.Errorf("status302: invalid stream URL %q", s.url))
	}
	switch s.encoding {
	case StreamEncodingJSON:
	case StreamEncodingMsgpack:
		s.url = withQuery(s.url, "encoding", string(s.encoding))
	default:
		errs = append(errs, fmt.Errorf("status302: unknown stream encoding %q", s.encoding))
	}
	apiKey := cmp.Or(opts.APIKey, c.apiKey)
	auth := opts.Auth
	if auth == "" {
		auth = StreamAuthTicket
		if apiKey != "" {
			auth = StreamAuthBearer
		}
	}
	switch auth {
	case StreamAuthBearer:
		if apiKey == "" {
			errs = append(errs, errors.New("status302: bearer stream auth needs an API key"))
		}
		s.header = http.Header{"Authorization": {"Bearer " + apiKey}}
	case StreamAuthTicket:
		s.getTicket = opts.GetTicket
		if s.getTicket == nil {
			s.getTicket = func(ctx context.Context) (string, error) {
				t, err := c.CreateStreamTicket(ctx)
				if err != nil {
					return "", err
				}
				return t.Ticket, nil
			}
		}
	default:
		errs = append(errs, fmt.Errorf("status302: unknown stream auth %q", auth))
	}
	if len(errs) > 0 {
		s.configErr = errs[0]
	}
	return s
}

// StreamURL derives the stream endpoint from an API origin:
// `https://api.status302.com` → `wss://api.status302.com/v1/stream` (`http` → `ws`).
func StreamURL(baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	scheme := "wss"
	if u.Scheme == "http" {
		scheme = "ws"
	}
	return (&url.URL{Scheme: scheme, User: u.User, Host: u.Host, Path: "/v1/stream"}).String()
}

func withQuery(u, key, value string) string {
	sep := "?"
	if strings.Contains(u, "?") {
		sep = "&"
	}
	return u + sep + key + "=" + url.QueryEscape(value)
}

// Subscribe subscribes to a topic; the first call connects. Subscribing to
// the same topic again replaces its subscription (and filter). Messages
// arrive through Next. After Close the returned Subscription's Ready fails
// with ErrStreamClosed.
func (s *Stream) Subscribe(topic string, opts SubscribeOptions) *Subscription {
	sub := &Subscription{stream: s, topic: topic, ready: make(chan struct{})}
	if len(opts.Filter) > 0 {
		filter, err := json.Marshal(opts.Filter)
		if err != nil {
			sub.settle(0, fmt.Errorf("status302: %s: encode filter: %w", topic, err))
			return sub
		}
		sub.filter = filter
	}
	if opts.Since != nil {
		sub.lastSeq, sub.hasSeq = *opts.Since, true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		sub.settle(0, ErrStreamClosed)
		return sub
	}
	s.subs[topic] = sub
	if !s.started {
		s.started = true
		s.startDispatchLocked()
		s.setStateLocked(StreamConnecting)
		go s.run()
	} else if s.conn != nil {
		s.sendSubscribeLocked(sub)
	}
	return sub
}

// Next returns the next message, waiting for one until ctx ends. Once the
// stream ended for good it returns the final error (a *StreamClosedError, or
// the invalid StreamOptions) once, after the messages queued before it, then
// ErrStreamClosed. After Close it returns ErrStreamClosed.
func (s *Stream) Next(ctx context.Context) (Message, error) {
	for {
		s.mu.Lock()
		switch {
		case s.userClosed:
			s.mu.Unlock()
			return nil, ErrStreamClosed
		case s.head < len(s.queue):
			m := s.queue[s.head]
			s.queue[s.head] = nil
			s.head++
			if s.head == len(s.queue) {
				s.queue, s.head = s.queue[:0], 0
			} else {
				signal(s.wake) // another caller may be waiting
			}
			s.mu.Unlock()
			return m, nil
		case s.closed:
			err := s.final
			s.final = nil
			s.mu.Unlock()
			if err == nil {
				err = ErrStreamClosed
			}
			return nil, err
		}
		s.mu.Unlock()
		select {
		case <-s.wake:
		case <-s.done:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Messages iterates over Next until the stream ends: it yields every message,
// then the final error or ctx.Err() (and stops); ErrStreamClosed ends it
// without being yielded.
func (s *Stream) Messages(ctx context.Context) iter.Seq2[Message, error] {
	return func(yield func(Message, error) bool) {
		for {
			m, err := s.Next(ctx)
			if errors.Is(err, ErrStreamClosed) || !yield(m, err) || err != nil {
				return
			}
		}
	}
}

// FetchSnapshot loads a resync snapshot through the REST client.
func (s *Stream) FetchSnapshot(ctx context.Context, m *ResyncMessage) (json.RawMessage, error) {
	return s.client.GetPath(ctx, m.Snapshot)
}

// State is the current connection state.
func (s *Stream) State() StreamState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// Close ends the stream for good: it closes the connection (1000 "client
// closed", waiting for the close handshake), drops undelivered messages and
// makes Next return ErrStreamClosed. Calling it again does nothing.
func (s *Stream) Close() error {
	s.mu.Lock()
	if s.userClosed {
		s.mu.Unlock()
		return nil
	}
	s.userClosed = true
	clear(s.queue)
	s.queue, s.head, s.final = nil, 0, nil
	if !s.closed {
		s.endLocked()
	}
	if !s.started {
		s.started = true
		s.startDispatchLocked()
	}
	c := s.conn
	s.conn = nil
	s.mu.Unlock()
	if c != nil {
		_ = c.ws.Close(websocket.StatusNormalClosure, "client closed")
	}
	s.cancel()
	return nil
}

// Subscription is one topic subscription; it survives reconnects.
type Subscription struct {
	stream *Stream
	topic  string
	filter json.RawMessage
	ready  chan struct{} // closed once confirmed or refused
	seq    int64         // the confirmed head seq, set before ready closes
	err    error         // the refusal, set before ready closes

	// Guarded by stream.mu.
	lastSeq int64
	hasSeq  bool
	resumed bool // resubscribed after a reconnect, not yet confirmed
	settled bool
}

// Topic is the subscribed topic.
func (sub *Subscription) Topic() string { return sub.topic }

// Ready waits until the server confirms the first subscribe and returns the
// topic's head seq. A refused subscription returns a *SubscriptionError, a
// stream that ended first ErrStreamClosed.
func (sub *Subscription) Ready(ctx context.Context) (int64, error) {
	select {
	case <-sub.ready:
		return sub.seq, sub.err
	default:
	}
	select {
	case <-sub.ready:
		return sub.seq, sub.err
	case <-sub.stream.done:
		select {
		case <-sub.ready:
			return sub.seq, sub.err
		default:
			return 0, ErrStreamClosed
		}
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// LastSeq is the seq of the last delivered message (or the confirmed head
// seq before any); resubscribes resume after it. False before the first
// confirmation and after a resync.
func (sub *Subscription) LastSeq() (int64, bool) {
	sub.stream.mu.Lock()
	defer sub.stream.mu.Unlock()
	return sub.lastSeq, sub.hasSeq
}

// Unsubscribe stops the subscription; the server is told when connected.
func (sub *Subscription) Unsubscribe() {
	s := sub.stream
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subs[sub.topic] != sub {
		return
	}
	delete(s.subs, sub.topic)
	s.sendLocked(clientMessage{Op: "unsubscribe", ID: s.idLocked(), Topic: sub.topic})
}

// settle resolves Ready once (under stream.mu once the subscription is registered).
func (sub *Subscription) settle(seq int64, err error) {
	if sub.settled {
		return
	}
	sub.settled = true
	sub.seq, sub.err = seq, err
	close(sub.ready)
}

// ---- connection loop ------------------------------------------------------

func (s *Stream) run() {
	if s.configErr != nil {
		s.fail(s.configErr)
		return
	}
	attempt := 0
	for {
		delay, err := s.connect(&attempt)
		if err != nil {
			s.fail(err)
			return
		}
		s.mu.Lock()
		closed := s.closed
		s.setStateLocked(StreamReconnecting)
		s.mu.Unlock()
		if closed || s.ctx.Err() != nil {
			return
		}
		attempt++
		if sleepCtx(s.ctx, delay) != nil {
			return
		}
	}
}

// connect opens one connection and serves it until it ends. It returns the
// delay before the next attempt, or the error that ends the stream.
func (s *Stream) connect(attempt *int) (time.Duration, error) {
	u := s.url
	if s.getTicket != nil {
		ticket, err := s.getTicket(s.ctx)
		if err != nil {
			if s.ctx.Err() != nil {
				return 0, nil //nolint:nilerr // closed while fetching a ticket: the stream ends without an error
			}
			// A refused credential is final; anything else is retried.
			if apiErr, ok := errors.AsType[*APIError](err); ok && (apiErr.StatusCode == http.StatusUnauthorized ||
				apiErr.StatusCode == http.StatusPaymentRequired || apiErr.StatusCode == http.StatusForbidden) {
				return 0, &StreamClosedError{Code: CloseUnauthorized, Reason: "ticket refused", err: err}
			}
			s.report(err)
			return s.backoff(*attempt), nil
		}
		u = withQuery(u, "ticket", ticket)
	}
	dialCtx, cancel := context.WithTimeout(s.ctx, s.idleTimeout)
	ws, _, err := websocket.Dial(dialCtx, u, &websocket.DialOptions{HTTPClient: s.httpClient, HTTPHeader: s.header}) //nolint:bodyclose // coder/websocket closes the handshake response body
	cancel()
	if err != nil {
		if s.ctx.Err() != nil {
			return 0, nil //nolint:nilerr // closed while dialing: the stream ends without an error
		}
		return s.backoff(*attempt), nil
	}
	// Server messages are trusted and topic snapshots have no fixed bound.
	ws.SetReadLimit(-1)
	c := &streamConn{ws: ws, wake: make(chan struct{}, 1)}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		_ = ws.Close(websocket.StatusNormalClosure, "client closed")
		return 0, nil
	}
	s.conn = c
	*attempt = 0
	resumed := s.connectedOnce
	s.connectedOnce = true
	s.setStateLocked(StreamOpen)
	for _, topic := range slices.Sorted(maps.Keys(s.subs)) {
		sub := s.subs[topic]
		sub.resumed = resumed
		s.sendSubscribeLocked(sub)
	}
	s.mu.Unlock()

	writeCtx, stopWriter := context.WithCancel(s.ctx)
	go s.write(writeCtx, c)
	err = s.read(c)
	stopWriter()
	_ = ws.CloseNow()

	s.mu.Lock()
	if s.conn == c {
		s.conn = nil
	}
	clear(s.pending)
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return 0, nil
	}
	ce, ok := errors.AsType[websocket.CloseError](err)
	if !ok {
		return s.backoff(*attempt), nil
	}
	switch ce.Code {
	case CloseUnauthorized, CloseForbidden:
		return 0, &StreamClosedError{Code: int(ce.Code), Reason: ce.Reason}
	case CloseReconnect:
		return 0, nil
	case CloseTooManyConnections:
		s.report(&StreamClosedError{Code: int(ce.Code), Reason: ce.Reason})
	}
	return s.backoff(*attempt), nil
}

// read handles the connection's messages until it ends; nothing arriving for
// IdleTimeout closes it with 4000 "idle".
func (s *Stream) read(c *streamConn) error {
	idle := time.AfterFunc(s.idleTimeout, func() { _ = c.ws.Close(closeIdle, "idle") })
	defer idle.Stop()
	for {
		typ, data, err := c.ws.Read(s.ctx)
		if err != nil {
			return err
		}
		idle.Reset(s.idleTimeout)
		s.handle(c, typ, data)
	}
}

// write sends the connection's queued client messages, in order.
func (s *Stream) write(ctx context.Context, c *streamConn) {
	var batch [][]byte
	for {
		select {
		case <-c.wake:
		case <-ctx.Done():
			return
		}
		s.mu.Lock()
		batch, c.out = c.out, batch[:0]
		s.mu.Unlock()
		for _, b := range batch {
			if err := c.ws.Write(ctx, websocket.MessageText, b); err != nil {
				_ = c.ws.CloseNow() // the read loop reconnects
				return
			}
		}
		clear(batch)
	}
}

func (s *Stream) handle(c *streamConn, typ websocket.MessageType, data []byte) {
	// Text frames are always JSON; binary frames are MessagePack with the msgpack encoding.
	if typ == websocket.MessageBinary && s.encoding == StreamEncodingMsgpack {
		var err error
		if data, err = msgpackToJSON(data); err != nil {
			s.report(fmt.Errorf("status302: malformed stream message: %w", err))
			return
		}
	}
	var m serverMessage
	if err := json.Unmarshal(data, &m); err != nil {
		s.report(fmt.Errorf("status302: malformed stream message: %w", err))
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	switch m.Type {
	case "subscribed":
		sub, ok := s.pending[m.ID]
		delete(s.pending, m.ID)
		if !ok {
			sub = s.subs[m.Topic]
		}
		if sub == nil || s.subs[sub.topic] != sub {
			return // replaced or unsubscribed meanwhile
		}
		if !sub.hasSeq {
			sub.lastSeq, sub.hasSeq = m.Seq, true
		}
		sub.settle(m.Seq, nil)
		if sub.resumed && coalescedTopic.MatchString(sub.topic) && m.Snapshot != "" {
			s.pushLocked(&ResyncMessage{Topic: sub.topic, Reason: "reconnect", Snapshot: m.Snapshot})
		}
		sub.resumed = false
	case "data":
		sub := s.subs[m.Topic]
		if sub == nil || (sub.hasSeq && m.Seq <= sub.lastSeq) {
			return // at-least-once delivery: drop duplicates
		}
		sub.lastSeq, sub.hasSeq = m.Seq, true
		if m.Data == nil {
			m.Data = jsonNull
		}
		s.pushLocked(&DataMessage{Topic: m.Topic, Seq: m.Seq, TS: m.TS, Data: m.Data})
	case "resync":
		sub := s.subs[m.Topic]
		if sub == nil {
			return
		}
		// The server continues after its current head, which may be below the
		// last seq seen (its counter restarted): deduplication starts over.
		sub.hasSeq = false
		s.pushLocked(&ResyncMessage{Topic: m.Topic, Reason: m.Reason, Snapshot: m.Snapshot})
	case "error":
		sub, ok := s.pending[m.ID]
		if !ok {
			s.reportLocked(fmt.Errorf("status302: stream error: %s: %s", m.Code, m.Message))
			return
		}
		delete(s.pending, m.ID)
		if s.subs[sub.topic] != sub {
			return // replaced or unsubscribed meanwhile
		}
		if m.Code == "unavailable" || m.Code == "rate_limited" {
			// Transient: retry this subscribe shortly on the same connection.
			time.AfterFunc(subscribeRetryDelay, func() {
				s.mu.Lock()
				defer s.mu.Unlock()
				if !s.closed && s.conn == c && s.subs[sub.topic] == sub {
					s.sendSubscribeLocked(sub)
				}
			})
			return
		}
		delete(s.subs, sub.topic)
		sub.settle(0, &SubscriptionError{Topic: sub.topic, Code: m.Code, Message: m.Message})
		s.pushLocked(&ErrorMessage{Topic: sub.topic, Code: m.Code, Message: m.Message})
	}
}

func (s *Stream) sendSubscribeLocked(sub *Subscription) {
	if s.conn == nil {
		return
	}
	msg := clientMessage{Op: "subscribe", ID: s.idLocked(), Topic: sub.topic, Filter: sub.filter}
	if sub.hasSeq && !coalescedTopic.MatchString(sub.topic) {
		msg.Since = new(sub.lastSeq)
	}
	s.pending[msg.ID] = sub
	s.sendLocked(msg)
}

// sendLocked queues a client message on the current connection, if any.
func (s *Stream) sendLocked(msg clientMessage) {
	if s.conn == nil {
		return
	}
	b, _ := json.Marshal(msg) // strings, an int and pre-encoded JSON: cannot fail
	s.conn.out = append(s.conn.out, b)
	signal(s.conn.wake)
}

func (s *Stream) idLocked() string {
	s.nextID++
	return "r" + strconv.Itoa(s.nextID)
}

func (s *Stream) pushLocked(m Message) {
	// Reuse the consumed front once it is at least half of a full queue.
	if s.head > 0 && s.head >= len(s.queue)/2 && len(s.queue) == cap(s.queue) {
		n := copy(s.queue, s.queue[s.head:])
		clear(s.queue[n:])
		s.queue, s.head = s.queue[:n], 0
	}
	s.queue = append(s.queue, m)
	signal(s.wake)
}

// backoff is exponential with full jitter: random in [0, min(MaxReconnectDelay, ReconnectDelay · 2^attempt)].
func (s *Stream) backoff(attempt int) time.Duration {
	ceiling := s.baseDelay
	for i := 0; i < attempt && ceiling < s.maxDelay; i++ {
		ceiling *= 2
	}
	ceiling = min(ceiling, s.maxDelay)
	if ceiling <= 0 {
		return 0
	}
	return time.Duration(rand.Int64N(int64(ceiling) + 1))
}

// fail ends the stream with a final error for Next.
func (s *Stream) fail(err error) {
	s.mu.Lock()
	if !s.closed {
		s.final = err
		s.endLocked()
	}
	s.mu.Unlock()
	s.cancel()
}

func (s *Stream) endLocked() {
	s.setStateLocked(StreamClosed)
	s.closed = true
	clear(s.pending)
	close(s.done)
}

// ---- callbacks ------------------------------------------------------------

func (s *Stream) setStateLocked(state StreamState) {
	if s.state == state || s.state == StreamClosed {
		return
	}
	s.state = state
	if s.onState != nil {
		s.events = append(s.events, streamEvent{state: state})
		signal(s.eventWake)
	}
}

func (s *Stream) report(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reportLocked(err)
}

func (s *Stream) reportLocked(err error) {
	if s.onError == nil || s.closed {
		return
	}
	s.events = append(s.events, streamEvent{err: err})
	signal(s.eventWake)
}

func (s *Stream) startDispatchLocked() {
	if s.onState != nil || s.onError != nil {
		go s.dispatch()
	}
}

// dispatch runs the callbacks in order, outside the lock, so they may call
// back into the stream; it drains the queue and returns once the stream ended.
func (s *Stream) dispatch() {
	var batch []streamEvent
	for {
		s.mu.Lock()
		batch, s.events = s.events, batch[:0]
		closed := s.closed
		s.mu.Unlock()
		for _, ev := range batch {
			if ev.err != nil {
				s.onError(ev.err)
			} else {
				s.onState(ev.state)
			}
		}
		if len(batch) == 0 {
			if closed {
				return
			}
			select {
			case <-s.eventWake:
			case <-s.done:
			}
		}
		clear(batch)
	}
}

func signal(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// ---- msgpack --------------------------------------------------------------

// msgpackToJSON transcodes one MessagePack value into JSON, keeping the map
// key order, so binary frames are handled exactly like text frames.
func msgpackToJSON(data []byte) ([]byte, error) {
	dec := msgpack.GetDecoder()
	defer msgpack.PutDecoder(dec)
	dec.Reset(bytes.NewReader(data))
	return appendMsgpackJSON(make([]byte, 0, len(data)+len(data)/2), dec, 0)
}

func appendMsgpackJSON(buf []byte, dec *msgpack.Decoder, depth int) ([]byte, error) {
	if depth > maxMsgpackDepth {
		return nil, errors.New("msgpack: nesting too deep")
	}
	c, err := dec.PeekCode()
	if err != nil {
		return nil, err
	}
	switch {
	case msgpcode.IsFixedMap(c) || c == msgpcode.Map16 || c == msgpcode.Map32:
		n, err := dec.DecodeMapLen()
		if err != nil {
			return nil, err
		}
		buf = append(buf, '{')
		for i := range n {
			if i > 0 {
				buf = append(buf, ',')
			}
			key, err := dec.DecodeString()
			if err != nil {
				return nil, err
			}
			if buf, err = appendJSON(buf, key); err != nil {
				return nil, err
			}
			buf = append(buf, ':')
			if buf, err = appendMsgpackJSON(buf, dec, depth+1); err != nil {
				return nil, err
			}
		}
		return append(buf, '}'), nil
	case msgpcode.IsFixedArray(c) || c == msgpcode.Array16 || c == msgpcode.Array32:
		n, err := dec.DecodeArrayLen()
		if err != nil {
			return nil, err
		}
		buf = append(buf, '[')
		for i := range n {
			if i > 0 {
				buf = append(buf, ',')
			}
			if buf, err = appendMsgpackJSON(buf, dec, depth+1); err != nil {
				return nil, err
			}
		}
		return append(buf, ']'), nil
	}
	// Scalars: ints become int64/uint64, floats float64, nil null.
	v, err := dec.DecodeInterfaceLoose()
	if err != nil {
		return nil, err
	}
	return appendJSON(buf, v)
}

func appendJSON(buf []byte, v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return append(buf, b...), nil
}
