package status302

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the production API origin.
const DefaultBaseURL = "https://api.status302.com"

const (
	baseDelay = 500 * time.Millisecond
	maxDelay  = 30 * time.Second
)

// Client is the Status302 REST client. Every operation of the public API is a
// method (operations.go), typed from the public OpenAPI spec:
//
//	client := status302.NewClient(os.Getenv("STATUS302_API_KEY"))
//	page, err := client.ListEvents(ctx, &status302.ListEventsParams{Limit: new(10)})
//
// A Client is safe for concurrent use.
type Client struct {
	apiKey        string
	baseURL       string
	httpClient    *http.Client
	maxRetries    int
	maxRetryAfter time.Duration
	timeout       time.Duration
	etagCache     ETagCache
	userAgent     string
	sleep         func(ctx context.Context, d time.Duration) error
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets the API origin (default DefaultBaseURL).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithHTTPClient sets the HTTP client (default http.DefaultClient). Its own
// timeout, if any, applies on top of WithTimeout.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithMaxRetries sets how often a request is retried after a 429, a 5xx
// (idempotent methods) or a connection error. Default 3.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.maxRetries = max(0, n) }
}

// WithMaxRetryAfter sets the longest Retry-After the client waits out before
// returning the 429 (or 503) error. Default 60 s.
func WithMaxRetryAfter(d time.Duration) Option {
	return func(c *Client) { c.maxRetryAfter = d }
}

// WithTimeout sets the per-attempt timeout. Default 30 s; 0 disables it.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.timeout = d }
}

// WithETagCache makes the client revalidate GET responses that carried an
// ETag: the next GET of the same URL sends If-None-Match and a 304 returns
// the stored body, so unchanged data costs no transfer. E.g.
// NewMemoryETagCache(1000). Off by default.
func WithETagCache(cache ETagCache) Option {
	return func(c *Client) { c.etagCache = cache }
}

// NewClient returns a client authenticating with apiKey
// (`tx_live_<keyId>_<secret>`, sent as `Authorization: Bearer`). An empty key
// sends no credentials (e.g. when the HTTP client adds them itself).
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:        apiKey,
		baseURL:       DefaultBaseURL,
		httpClient:    http.DefaultClient,
		maxRetries:    3,
		maxRetryAfter: 60 * time.Second,
		timeout:       30 * time.Second,
		userAgent:     "status302-go/" + Version,
		sleep:         sleepCtx,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// BaseURL is the API origin the client sends requests to.
func (c *Client) BaseURL() string { return c.baseURL }

// Response describes the HTTP response of a call: pass WithResponseInto to receive it.
type Response struct {
	// StatusCode is the HTTP status; 304 when the body was revalidated from the client's ETag cache.
	StatusCode int
	Header     http.Header
	// ETag is the response's ETag header ("" when absent): pass it to WithIfNoneMatch to revalidate.
	ETag string
	// RateLimit is the parsed RateLimit / RateLimit-Policy state (nil on unmetered requests).
	RateLimit *RateLimitInfo
}

// RequestOption adjusts a single call.
type RequestOption func(*requestConfig)

type requestConfig struct {
	header      http.Header
	ifNoneMatch string
	response    *Response
}

// WithHeader adds a request header to the call.
func WithHeader(key, value string) RequestOption {
	return func(r *requestConfig) {
		if r.header == nil {
			r.header = http.Header{}
		}
		r.header.Add(key, value)
	}
}

// WithIfNoneMatch makes the call conditional: while the representation still
// has that ETag the API answers 304 and the call fails with a
// *NotModifiedError (errors.Is(err, ErrNotModified)). An empty etag sends no
// condition, so WithIfNoneMatch(previous.ETag) works on the first call too.
func WithIfNoneMatch(etag string) RequestOption {
	return func(r *requestConfig) { r.ifNoneMatch = etag }
}

// WithResponseInto stores the call's status, headers, ETag and rate limit state in resp.
func WithResponseInto(resp *Response) RequestOption {
	return func(r *requestConfig) { r.response = resp }
}

// GetPath GETs a REST path returned by the API itself (a stream snapshot, a
// resync target), including its query string, and returns the JSON body.
func (c *Client) GetPath(ctx context.Context, pathAndQuery string, opts ...RequestOption) (json.RawMessage, error) {
	if !strings.HasPrefix(pathAndQuery, "/") {
		return nil, fmt.Errorf("status302: GetPath expects a path starting with /, got %q", pathAndQuery)
	}
	var out json.RawMessage
	if err := c.send(ctx, http.MethodGet, c.baseURL+pathAndQuery, nil, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}

// jsonBody is the body argument of an operation: a nil pointer sends no body.
func jsonBody[T any](v *T) any {
	if v == nil {
		return nil
	}
	return v
}

// do sends a request to path (already escaped) with the query; out is nil
// (no body), *[]byte (raw body) or a pointer to decode JSON into.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any, opts []RequestOption) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return c.send(ctx, method, u, body, out, opts)
}

func (c *Client) send(ctx context.Context, method, u string, body, out any, opts []RequestOption) error {
	var cfg requestConfig
	for _, o := range opts {
		o(&cfg)
	}
	var payload []byte
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("status302: encode request body: %w", err)
		}
		payload = data
	}

	header := http.Header{}
	header.Set("Accept", "application/json, application/problem+json")
	header.Set("User-Agent", c.userAgent)
	if c.apiKey != "" {
		header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if payload != nil {
		header.Set("Content-Type", "application/json")
	}
	for k, vs := range cfg.header {
		header[k] = vs
	}
	if cfg.ifNoneMatch != "" {
		header.Set("If-None-Match", cfg.ifNoneMatch)
	}
	// The cache revalidates only requests without a caller-chosen condition.
	var cache ETagCache
	if method == http.MethodGet && header.Get("If-None-Match") == "" {
		cache = c.etagCache
	}
	var cached *CachedResponse
	if cache != nil {
		if entry, ok := cache.Get(ctx, u); ok {
			cached = &entry
			header.Set("If-None-Match", entry.ETag)
		}
	}
	idempotent := isIdempotent(method)

	for attempt := 0; ; attempt++ {
		res, data, err := c.attempt(ctx, method, u, header, payload)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if idempotent && attempt < c.maxRetries {
				if err := c.sleep(ctx, backoff(attempt)); err != nil {
					return err
				}
				continue
			}
			return &ConnectionError{Method: method, URL: u, Err: err}
		}

		rateLimit := ParseRateLimit(res.Header)
		etag := res.Header.Get("ETag")
		meta := Response{StatusCode: res.StatusCode, Header: res.Header, ETag: etag, RateLimit: rateLimit}
		if cfg.response != nil {
			*cfg.response = meta
		}
		switch {
		case res.StatusCode == http.StatusNotModified:
			if cached == nil {
				return &NotModifiedError{ETag: etag, Header: res.Header, RateLimit: rateLimit}
			}
			if cfg.response != nil && cfg.response.ETag == "" {
				cfg.response.ETag = cached.ETag
			}
			return decodeBody(cached.Body, cached.ContentType, out)
		case res.StatusCode >= 200 && res.StatusCode < 300:
			contentType := res.Header.Get("Content-Type")
			if cache != nil && etag != "" && res.StatusCode == http.StatusOK {
				cache.Set(ctx, u, CachedResponse{ETag: etag, Body: data, ContentType: contentType})
			}
			if res.StatusCode == http.StatusNoContent {
				return nil
			}
			return decodeBody(data, contentType, out)
		}

		retryAfter, hasRetryAfter := ParseRetryAfter(res.Header.Get("Retry-After"), time.Now())
		apiErr := newAPIError(res, data, rateLimit, retryAfter, hasRetryAfter)
		delay, retry := c.retryDelay(res.StatusCode, idempotent, attempt, retryAfter, hasRetryAfter)
		if !retry {
			return apiErr
		}
		if err := c.sleep(ctx, delay); err != nil {
			return err
		}
	}
}

// attempt performs one HTTP exchange and reads the whole body.
func (c *Client) attempt(ctx context.Context, method, u string, header http.Header, payload []byte) (*httpReply, []byte, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header = header.Clone()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return &httpReply{StatusCode: res.StatusCode, Header: res.Header}, data, nil
}

// httpReply is the status and headers of a response whose body was read.
type httpReply struct {
	StatusCode int
	Header     http.Header
}

// retryDelay is the wait before retrying, or false when the error is final.
func (c *Client) retryDelay(status int, idempotent bool, attempt int, retryAfter time.Duration, hasRetryAfter bool) (time.Duration, bool) {
	if attempt >= c.maxRetries {
		return 0, false
	}
	// 429 means the request was refused before it ran, so any method may retry it.
	retryable := status == http.StatusTooManyRequests || (status >= 500 && idempotent) || (status == http.StatusServiceUnavailable && hasRetryAfter)
	if !retryable {
		return 0, false
	}
	if hasRetryAfter {
		return retryAfter, retryAfter <= c.maxRetryAfter
	}
	return backoff(attempt), true
}

func decodeBody(data []byte, contentType string, out any) error {
	switch o := out.(type) {
	case nil:
		return nil
	case *[]byte:
		*o = data
		return nil
	case *json.RawMessage:
		*o = append(json.RawMessage(nil), data...)
		return nil
	}
	if !strings.Contains(contentType, "json") && contentType != "" {
		return fmt.Errorf("status302: unexpected response content type %q", contentType)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("status302: decode response: %w", err)
	}
	return nil
}

func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

// backoff is exponential with full jitter: random in [0, min(30 s, 0.5 s · 2^attempt)].
func backoff(attempt int) time.Duration {
	ceiling := maxDelay
	if attempt < 16 {
		ceiling = min(maxDelay, baseDelay<<attempt)
	}
	return time.Duration(rand.Int64N(int64(ceiling) + 1))
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
