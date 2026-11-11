package status302

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

const streamTestKey = "tx_live_abcdefghijkl_secret"

// fakeStreamServer is a scripted stream server: it records connections and
// client messages.
type fakeStreamServer struct {
	t   *testing.T
	url string

	mu      sync.Mutex
	changed chan struct{} // closed and replaced on every change
	conns   []*fakeStreamConn
	binary  int // client frames that arrived as binary (the client only ever sends text)
}

type fakeStreamConn struct {
	ws     *websocket.Conn
	path   string
	query  string
	header http.Header
	// Guarded by fakeStreamServer.mu.
	received []map[string]any
	closeErr error
}

func newFakeStreamServer(t *testing.T) *fakeStreamServer {
	s := &fakeStreamServer{t: t, changed: make(chan struct{})}
	srv := httptest.NewServer(http.HandlerFunc(s.serve))
	t.Cleanup(srv.Close)
	s.url = "ws" + strings.TrimPrefix(srv.URL, "http") + "/v1/stream"
	return s
}

func (s *fakeStreamServer) serve(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	c := &fakeStreamConn{ws: ws, path: r.URL.Path, query: r.URL.RawQuery, header: r.Header.Clone()}
	s.update(func() { s.conns = append(s.conns, c) })
	for {
		typ, data, err := ws.Read(context.Background())
		if err != nil {
			s.update(func() { c.closeErr = err })
			return
		}
		var msg map[string]any
		_ = json.Unmarshal(data, &msg)
		s.update(func() {
			if typ == websocket.MessageBinary {
				s.binary++
			}
			c.received = append(c.received, msg)
		})
	}
}

func (s *fakeStreamServer) update(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f()
	close(s.changed)
	s.changed = make(chan struct{})
}

// until waits up to 3 s for pred to hold.
func (s *fakeStreamServer) until(pred func() bool) {
	s.t.Helper()
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	for {
		s.mu.Lock()
		changed := s.changed
		s.mu.Unlock()
		if pred() {
			return
		}
		select {
		case <-changed:
		case <-deadline.C:
			s.t.Fatal("timed out waiting for the stream server")
		}
	}
}

func (s *fakeStreamServer) connCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.conns)
}

func (s *fakeStreamServer) conn(i int) *fakeStreamConn {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conns[i]
}

func (s *fakeStreamServer) binaryFrames() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.binary
}

func (s *fakeStreamServer) closeErr(i int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conns[i].closeErr
}

// subscribes returns the subscribe requests received on connection i.
func (s *fakeStreamServer) subscribes(i int) []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.conns) <= i {
		return nil
	}
	var out []map[string]any
	for _, m := range s.conns[i].received {
		if m["op"] == "subscribe" {
			out = append(out, m)
		}
	}
	return out
}

func (s *fakeStreamServer) send(i int, msg map[string]any) {
	s.t.Helper()
	b, err := json.Marshal(msg)
	if err != nil {
		s.t.Fatal(err)
	}
	s.write(i, websocket.MessageText, b)
}

// sendPacked sends one binary frame holding one MessagePack map (the msgpack
// encoding), keys sorted like encoding/json sorts them.
func (s *fakeStreamServer) sendPacked(i int, msg map[string]any) {
	s.t.Helper()
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetSortMapKeys(true)
	if err := enc.Encode(msg); err != nil {
		s.t.Fatal(err)
	}
	s.write(i, websocket.MessageBinary, buf.Bytes())
}

func (s *fakeStreamServer) write(i int, typ websocket.MessageType, b []byte) {
	s.t.Helper()
	if err := s.conn(i).ws.Write(s.t.Context(), typ, b); err != nil {
		s.t.Fatal(err)
	}
}

func (s *fakeStreamServer) close(i int, code websocket.StatusCode, reason string) {
	_ = s.conn(i).ws.Close(code, reason)
}

func testCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func nextMessage(t *testing.T, stream *Stream) Message {
	t.Helper()
	m, err := stream.Next(testCtx(t))
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	return m
}

func nextData(t *testing.T, stream *Stream) *DataMessage {
	t.Helper()
	m := nextMessage(t, stream)
	d, ok := m.(*DataMessage)
	if !ok {
		t.Fatalf("got %#v, want *DataMessage", m)
	}
	return d
}

func TestStreamBearerResumeAndDedupe(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient(streamTestKey).Stream(StreamOptions{URL: s.url, ReconnectDelay: 10 * time.Millisecond})
	defer stream.Close()

	sub := stream.Subscribe("sales", SubscribeOptions{Filter: map[string]any{"country": []string{"US"}}})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	if got := s.conn(0).header.Get("Authorization"); got != "Bearer "+streamTestKey {
		t.Fatalf("Authorization = %q", got)
	}
	first := s.subscribes(0)[0]
	wantFilter := map[string]any{"country": []any{"US"}}
	if first["topic"] != "sales" || !reflect.DeepEqual(first["filter"], wantFilter) {
		t.Fatalf("subscribe = %v", first)
	}
	if _, ok := first["since"]; ok {
		t.Fatalf("first subscribe has since: %v", first)
	}

	s.send(0, map[string]any{"type": "subscribed", "id": first["id"], "topic": "sales", "seq": 10})
	if seq, err := sub.Ready(testCtx(t)); err != nil || seq != 10 {
		t.Fatalf("Ready = %d, %v", seq, err)
	}
	for _, seq := range []int{11, 12, 12, 11, 13} {
		s.send(0, map[string]any{"type": "data", "topic": "sales", "seq": seq, "ts": "2026-10-28T00:00:00Z", "data": map[string]any{"seq": seq}})
	}
	for _, want := range []int64{11, 12, 13} {
		if m := nextData(t, stream); m.Seq != want {
			t.Fatalf("seq = %d, want %d", m.Seq, want)
		}
	}

	s.close(0, CloseSlowConsumer, "slow_consumer")
	s.until(func() bool { return len(s.subscribes(1)) == 1 })
	resumed := s.subscribes(1)[0]
	if resumed["since"] != float64(13) || !reflect.DeepEqual(resumed["filter"], wantFilter) {
		t.Fatalf("resubscribe = %v", resumed)
	}
	s.send(1, map[string]any{"type": "subscribed", "id": resumed["id"], "topic": "sales", "seq": 15})
	for _, seq := range []int{13, 14, 15} {
		s.send(1, map[string]any{"type": "data", "topic": "sales", "seq": seq, "ts": "2026-10-28T00:00:01Z", "data": map[string]any{}})
	}
	for _, want := range []int64{14, 15} {
		if m := nextData(t, stream); m.Seq != want {
			t.Fatalf("seq = %d, want %d", m.Seq, want)
		}
	}
	if seq, ok := sub.LastSeq(); !ok || seq != 15 {
		t.Fatalf("LastSeq = %d, %v", seq, ok)
	}
}

func TestStreamTicketAuthFetchesATicketPerConnection(t *testing.T) {
	s := newFakeStreamServer(t)
	var tickets atomic.Int32
	stream := NewClient("").Stream(StreamOptions{
		URL:            s.url,
		ReconnectDelay: 10 * time.Millisecond,
		GetTicket: func(context.Context) (string, error) {
			return fmt.Sprintf("t%d", tickets.Add(1)), nil
		},
	})
	defer stream.Close()

	stream.Subscribe("onsales", SubscribeOptions{})
	s.until(func() bool { return s.connCount() == 1 })
	c := s.conn(0)
	if c.path != "/v1/stream" || c.query != "ticket=t1" {
		t.Fatalf("request = %s?%s", c.path, c.query)
	}
	if got := c.header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization = %q", got)
	}
	s.close(0, CloseReconnect, "reconnect")
	s.until(func() bool { return s.connCount() == 2 })
	if c := s.conn(1); c.query != "ticket=t2" {
		t.Fatalf("second request query = %q", c.query)
	}
}

func TestStreamTicketErrorIsReportedAndRetried(t *testing.T) {
	s := newFakeStreamServer(t)
	var calls atomic.Int32
	rest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/stream/ticket" {
			http.NotFound(w, r)
			return
		}
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ticket":"fresh ticket","expiresAt":"2026-11-11T00:01:00Z"}`))
	}))
	defer rest.Close()
	errs := make(chan error, 4)
	stream := NewClient("", WithBaseURL(rest.URL), WithMaxRetries(0)).Stream(StreamOptions{
		URL:            s.url,
		ReconnectDelay: 10 * time.Millisecond,
		OnError:        func(err error) { errs <- err },
	})
	defer stream.Close()

	stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return s.connCount() == 1 })
	if c := s.conn(0); c.query != "ticket=fresh+ticket" {
		t.Fatalf("query = %q", c.query)
	}
	select {
	case err := <-errs:
		if apiErr, ok := errors.AsType[*APIError](err); !ok || apiErr.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("OnError got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("OnError was not called")
	}
}

func TestStreamTicketRefusedIsFinal(t *testing.T) {
	s := newFakeStreamServer(t)
	rest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"type":"about:blank","title":"Unauthorized","status":401}`))
	}))
	defer rest.Close()
	stream := NewClient("", WithBaseURL(rest.URL)).Stream(StreamOptions{URL: s.url, ReconnectDelay: 10 * time.Millisecond})
	defer stream.Close()

	sub := stream.Subscribe("sales", SubscribeOptions{})
	_, err := stream.Next(testCtx(t))
	closed, ok := errors.AsType[*StreamClosedError](err)
	if !ok || closed.Code != CloseUnauthorized || closed.Reason != "ticket refused" {
		t.Fatalf("Next error = %v", err)
	}
	if apiErr, ok := errors.AsType[*APIError](err); !ok || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("cause = %v", err)
	}
	if _, err := stream.Next(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("second Next error = %v", err)
	}
	if _, err := sub.Ready(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("Ready error = %v", err)
	}
	if n := s.connCount(); n != 0 {
		t.Fatalf("%d connections", n)
	}
}

func TestStreamResyncAndSnapshot(t *testing.T) {
	s := newFakeStreamServer(t)
	var mu sync.Mutex
	var fetched []string
	rest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		fetched = append(fetched, r.Method+" "+r.URL.RequestURI())
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[],"nextCursor":null}`))
	}))
	defer rest.Close()
	stream := NewClient(streamTestKey, WithBaseURL(rest.URL)).Stream(StreamOptions{URL: s.url})
	defer stream.Close()

	stream.Subscribe("event.46.sales", SubscribeOptions{Since: new(int64(3))})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	req := s.subscribes(0)[0]
	if req["since"] != float64(3) {
		t.Fatalf("subscribe = %v", req)
	}
	s.send(0, map[string]any{"type": "subscribed", "id": req["id"], "topic": "event.46.sales", "seq": 900})
	s.send(0, map[string]any{"type": "resync", "topic": "event.46.sales", "reason": "gap", "snapshot": "/v1/events/46/sales"})
	msg := nextMessage(t, stream)
	resync, ok := msg.(*ResyncMessage)
	if want := (ResyncMessage{Topic: "event.46.sales", Reason: "gap", Snapshot: "/v1/events/46/sales"}); !ok || *resync != want {
		t.Fatalf("got %#v", msg)
	}
	body, err := stream.FetchSnapshot(testCtx(t), resync)
	if err != nil || string(body) != `{"items":[],"nextCursor":null}` {
		t.Fatalf("FetchSnapshot = %s, %v", body, err)
	}
	mu.Lock()
	if want := []string{"GET /v1/events/46/sales"}; !reflect.DeepEqual(fetched, want) {
		t.Fatalf("fetched %v", fetched)
	}
	mu.Unlock()
	// The server's counter restarted: seqs start over after the resync.
	s.send(0, map[string]any{"type": "data", "topic": "event.46.sales", "seq": 1, "ts": "2026-10-28T00:00:00Z", "data": map[string]any{"eventId": 46}})
	data := nextData(t, stream)
	var payload struct{ EventID int }
	if err := data.Decode(&payload); err != nil || data.Seq != 1 || payload.EventID != 46 {
		t.Fatalf("data = %+v (%s), %v", data, data.Data, err)
	}
}

func TestStreamCoalescedTopicReconnectResync(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient(streamTestKey).Stream(StreamOptions{URL: s.url, ReconnectDelay: 10 * time.Millisecond})
	defer stream.Close()

	stream.Subscribe("event.46.kpis", SubscribeOptions{})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	s.send(0, map[string]any{"type": "subscribed", "id": s.subscribes(0)[0]["id"], "topic": "event.46.kpis", "seq": 4, "snapshot": "/v1/events/46/metrics"})
	s.send(0, map[string]any{"type": "data", "topic": "event.46.kpis", "seq": 5, "ts": "2026-10-28T00:00:00Z", "data": map[string]any{"eventId": 46}})
	nextData(t, stream)
	s.close(0, CloseServiceRestart, "restart")
	s.until(func() bool { return len(s.subscribes(1)) == 1 })
	req := s.subscribes(1)[0]
	if _, ok := req["since"]; ok {
		t.Fatalf("coalesced resubscribe has since: %v", req)
	}
	s.send(1, map[string]any{"type": "subscribed", "id": req["id"], "topic": "event.46.kpis", "seq": 9, "snapshot": "/v1/events/46/metrics"})
	msg := nextMessage(t, stream)
	want := ResyncMessage{Topic: "event.46.kpis", Reason: "reconnect", Snapshot: "/v1/events/46/metrics"}
	if resync, ok := msg.(*ResyncMessage); !ok || *resync != want {
		t.Fatalf("got %#v", msg)
	}
}

func TestStreamUnauthorizedCloseIsFinal(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient("revoked").Stream(StreamOptions{URL: s.url, ReconnectDelay: 10 * time.Millisecond})
	defer stream.Close()

	stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	s.close(0, CloseUnauthorized, "unauthorized")
	_, err := stream.Next(testCtx(t))
	closed, ok := errors.AsType[*StreamClosedError](err)
	if !ok || closed.Code != CloseUnauthorized || closed.Reason != "unauthorized" {
		t.Fatalf("Next error = %v", err)
	}
	if got := err.Error(); got != "status302: stream closed (4401 unauthorized)" {
		t.Fatalf("Error() = %q", got)
	}
	if _, err := stream.Next(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("second Next error = %v", err)
	}
	if state := stream.State(); state != StreamClosed {
		t.Fatalf("state = %s", state)
	}
	time.Sleep(50 * time.Millisecond) // several reconnect delays
	if n := s.connCount(); n != 1 {
		t.Fatalf("%d connections", n)
	}
}

func TestStreamTooManyConnectionsIsReportedAndRetried(t *testing.T) {
	s := newFakeStreamServer(t)
	errs := make(chan error, 4)
	stream := NewClient(streamTestKey).Stream(StreamOptions{
		URL:            s.url,
		ReconnectDelay: 10 * time.Millisecond,
		OnError:        func(err error) { errs <- err },
	})
	defer stream.Close()

	stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return s.connCount() == 1 })
	s.close(0, CloseTooManyConnections, "too_many_connections")
	select {
	case err := <-errs:
		if closed, ok := errors.AsType[*StreamClosedError](err); !ok || closed.Code != CloseTooManyConnections {
			t.Fatalf("OnError got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("OnError was not called")
	}
	s.until(func() bool { return len(s.subscribes(1)) == 1 })
}

func TestStreamIdleTimeoutReconnects(t *testing.T) {
	s := newFakeStreamServer(t)
	states := make(chan StreamState, 16)
	stream := NewClient(streamTestKey).Stream(StreamOptions{
		URL:            s.url,
		ReconnectDelay: 10 * time.Millisecond,
		IdleTimeout:    100 * time.Millisecond,
		OnStateChange:  func(state StreamState) { states <- state },
	})

	stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return len(s.subscribes(1)) == 1 })
	if err := s.closeErrWait(0); websocket.CloseStatus(err) != closeIdle {
		t.Fatalf("first connection ended with %v", err)
	}
	stream.Close()
	// Later connections idle out too until Close: check the first cycle and the end.
	var got []StreamState
	timeout := time.After(3 * time.Second)
	for len(got) == 0 || got[len(got)-1] != StreamClosed {
		select {
		case state := <-states:
			got = append(got, state)
		case <-timeout:
			t.Fatalf("no %s state, got %v", StreamClosed, got)
		}
	}
	want := []StreamState{StreamConnecting, StreamOpen, StreamReconnecting, StreamOpen}
	if len(got) < len(want)+1 || !reflect.DeepEqual(got[:len(want)], want) {
		t.Fatalf("states = %v, want %v … %s", got, want, StreamClosed)
	}
}

// closeErrWait waits for connection i to end and returns its read error.
func (s *fakeStreamServer) closeErrWait(i int) error {
	s.t.Helper()
	s.until(func() bool { return s.closeErr(i) != nil })
	return s.closeErr(i)
}

func TestStreamRefusedSubscription(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient(streamTestKey).Stream(StreamOptions{URL: s.url})
	defer stream.Close()

	sub := stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	s.send(0, map[string]any{"type": "error", "id": s.subscribes(0)[0]["id"], "code": "forbidden_topic", "message": "requires plan Elite"})
	msg := nextMessage(t, stream)
	want := ErrorMessage{Topic: "sales", Code: "forbidden_topic", Message: "requires plan Elite"}
	if refused, ok := msg.(*ErrorMessage); !ok || *refused != want {
		t.Fatalf("got %#v", msg)
	}
	_, err := sub.Ready(testCtx(t))
	subErr, ok := errors.AsType[*SubscriptionError](err)
	if !ok || subErr.Topic != "sales" || subErr.Code != "forbidden_topic" {
		t.Fatalf("Ready error = %v", err)
	}
	if got := err.Error(); got != "status302: sales: forbidden_topic: requires plan Elite" {
		t.Fatalf("Error() = %q", got)
	}
}

// streamPayload is nested, with ints, floats, null, arrays, booleans and strings.
var streamPayload = map[string]any{
	"eventId": 46,
	"price":   12345,
	"fee":     12.5,
	"ratio":   -0.25,
	"section": nil,
	"seats":   []any{1, 2, 3},
	"tags":    []any{"vip", "aisle"},
	"resale":  true,
	"nested":  map[string]any{"rows": []any{map[string]any{"row": "A", "count": 2, "avg": 99.95}}, "empty": []any{}, "none": nil},
}

func TestStreamMsgpackEncoding(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient("").Stream(StreamOptions{
		URL:       s.url,
		Encoding:  StreamEncodingMsgpack,
		GetTicket: func(context.Context) (string, error) { return "t1", nil },
	})
	defer stream.Close()

	sub := stream.Subscribe("sales", SubscribeOptions{Filter: map[string]any{"country": []string{"US"}}})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	if c := s.conn(0); c.path != "/v1/stream" || c.query != "encoding=msgpack&ticket=t1" {
		t.Fatalf("request = %s?%s", c.path, c.query)
	}
	first := s.subscribes(0)[0]
	if first["topic"] != "sales" || !reflect.DeepEqual(first["filter"], map[string]any{"country": []any{"US"}}) {
		t.Fatalf("subscribe = %v", first)
	}
	if n := s.binaryFrames(); n != 0 { // the client's subscribe is a JSON text frame
		t.Fatalf("%d binary client frames", n)
	}

	s.sendPacked(0, map[string]any{"type": "subscribed", "id": first["id"], "topic": "sales", "seq": 10, "snapshot": "/v1/sales/feed"})
	if seq, err := sub.Ready(testCtx(t)); err != nil || seq != 10 {
		t.Fatalf("Ready = %d, %v", seq, err)
	}
	s.sendPacked(0, map[string]any{"type": "heartbeat", "ts": "2026-10-28T00:00:00Z"})
	for _, seq := range []int{11, 12, 11} {
		s.sendPacked(0, map[string]any{"type": "data", "topic": "sales", "seq": seq, "ts": "2026-10-28T00:00:00Z", "data": streamPayload})
	}
	// Defensive: a text frame is still parsed as JSON.
	s.send(0, map[string]any{"type": "data", "topic": "sales", "seq": 13, "ts": "2026-10-28T00:00:01Z", "data": streamPayload})
	got := []*DataMessage{nextData(t, stream), nextData(t, stream), nextData(t, stream)}
	for i, want := range []int64{11, 12, 13} {
		if got[i].Topic != "sales" || got[i].Seq != want {
			t.Fatalf("message %d = %+v", i, got[i])
		}
	}
	if got[0].TS != "2026-10-28T00:00:00Z" {
		t.Fatalf("ts = %q", got[0].TS)
	}
	// The decoded msgpack payload is the JSON path's, byte for byte (same key order).
	wantJSON, err := json.Marshal(streamPayload)
	if err != nil {
		t.Fatal(err)
	}
	for i, m := range got {
		if !bytes.Equal(m.Data, wantJSON) {
			t.Fatalf("message %d data = %s, want %s", i, m.Data, wantJSON)
		}
	}
	var decoded, want any
	if err := got[1].Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	_ = json.Unmarshal(wantJSON, &want)
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded %v, want %v", decoded, want)
	}
	if n := s.binaryFrames(); n != 0 {
		t.Fatalf("%d binary client frames", n)
	}
}

func TestStreamJSONEncodingAddsNoParameter(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient(streamTestKey).Stream(StreamOptions{URL: s.url})
	defer stream.Close()

	stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return s.connCount() == 1 })
	if c := s.conn(0); c.path != "/v1/stream" || c.query != "" {
		t.Fatalf("request = %s?%s", c.path, c.query)
	}
}

func TestStreamClose(t *testing.T) {
	s := newFakeStreamServer(t)
	stream := NewClient(streamTestKey).Stream(StreamOptions{URL: s.url})

	pending := stream.Subscribe("sales", SubscribeOptions{})
	s.until(func() bool { return len(s.subscribes(0)) == 1 })
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	err := s.closeErrWait(0)
	if ce, ok := errors.AsType[websocket.CloseError](err); !ok || ce.Code != websocket.StatusNormalClosure || ce.Reason != "client closed" {
		t.Fatalf("server saw %v", err)
	}
	if _, err := stream.Next(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("Next error = %v", err)
	}
	if _, err := pending.Ready(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("Ready of the unconfirmed subscription = %v", err)
	}
	if _, err := stream.Subscribe("onsales", SubscribeOptions{}).Ready(testCtx(t)); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("Ready after Close = %v", err)
	}
	if state := stream.State(); state != StreamClosed {
		t.Fatalf("state = %s", state)
	}
	for msg, err := range stream.Messages(testCtx(t)) {
		t.Fatalf("Messages yielded %v, %v after Close", msg, err)
	}
}

func TestStreamURL(t *testing.T) {
	for base, want := range map[string]string{
		"https://api.status302.com":     "wss://api.status302.com/v1/stream",
		"https://api.status302.com/":    "wss://api.status302.com/v1/stream",
		"http://localhost:8080":         "ws://localhost:8080/v1/stream",
		"http://localhost:8080/api?x=1": "ws://localhost:8080/v1/stream",
	} {
		if got := StreamURL(base); got != want {
			t.Errorf("StreamURL(%q) = %q, want %q", base, got, want)
		}
	}
}
