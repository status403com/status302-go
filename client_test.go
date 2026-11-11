package status302

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

const testKey = "tx_live_abcdefghijkl_secret"

// reply is one scripted answer of the fake transport: a response or a transport error.
type reply struct {
	status int
	header map[string]string
	body   any
	err    error
}

func problem(status int, header map[string]string, body map[string]any) reply {
	h := map[string]string{"Content-Type": "application/problem+json"}
	for k, v := range header {
		h[k] = v
	}
	body["status"] = status
	return reply{status: status, header: h, body: body}
}

func ok(body any, header map[string]string) reply {
	return reply{status: 200, header: header, body: body}
}

func notModified(etag string, header map[string]string) reply {
	h := map[string]string{"ETag": etag}
	for k, v := range header {
		h[k] = v
	}
	return reply{status: 304, header: h}
}

// recorder answers requests from a queue and records them.
type recorder struct {
	t        *testing.T
	replies  []reply
	requests []*http.Request
	bodies   [][]byte
}

func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	r.requests = append(r.requests, req)
	r.bodies = append(r.bodies, body)
	if len(r.replies) == 0 {
		r.t.Fatalf("unexpected request %s %s", req.Method, req.URL)
	}
	next := r.replies[0]
	r.replies = r.replies[1:]
	if next.err != nil {
		return nil, next.err
	}
	h := http.Header{}
	var data []byte
	if next.body != nil {
		data, _ = json.Marshal(next.body)
		h.Set("Content-Type", "application/json")
	}
	for k, v := range next.header {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: next.status, Header: h, Body: io.NopCloser(bytes.NewReader(data)), Request: req}, nil
}

func (r *recorder) urls() []string {
	var out []string
	for _, req := range r.requests {
		out = append(out, req.URL.String())
	}
	return out
}

func newTestClient(t *testing.T, replies []reply, opts ...Option) (*Client, *recorder, *[]time.Duration) {
	rec := &recorder{t: t, replies: replies}
	var sleeps []time.Duration
	c := NewClient(testKey, append([]Option{WithBaseURL("http://api.test/"), WithHTTPClient(&http.Client{Transport: rec})}, opts...)...)
	c.sleep = func(_ context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}
	return c, rec, &sleeps
}

func emptyPage() reply { return ok(map[string]any{"items": []any{}, "nextCursor": nil}, nil) }

func TestBearerKeyAndQueryEncoding(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{emptyPage()})
	page, err := c.ListEvents(t.Context(), &ListEventsParams{
		Q: new("taylor swift"), Country: new("US"), Limit: new(5), TmOnly: new(true),
		Status: []EventStatus{"onsale", "offsale"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 0 || page.NextCursor != nil {
		t.Fatalf("page %+v", page)
	}
	req := rec.requests[0]
	if got := req.Header.Get("Authorization"); got != "Bearer "+testKey {
		t.Errorf("authorization %q", got)
	}
	if got := req.Header.Get("User-Agent"); got != "status302-go/"+Version {
		t.Errorf("user-agent %q", got)
	}
	if got, want := req.URL.String(), "http://api.test/v1/events?country=US&limit=5&q=taylor+swift&status=onsale%2Coffsale&tmOnly=true"; got != want {
		t.Errorf("url %s, want %s", got, want)
	}
}

func TestPathParametersAndExplodedArrays(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{emptyPage(), ok(map[string]any{"id": "a b"}, nil)})
	if _, err := c.ListInventoryChanges(t.Context(), 46, &ListInventoryChangesParams{
		Section: []string{"101", "PIT B"}, Kinds: []InventoryChangeKind{"restock", "sale"},
	}); err != nil {
		t.Fatal(err)
	}
	u := rec.requests[0].URL
	if u.Path != "/v1/events/46/inventory/changes" {
		t.Errorf("path %s", u.Path)
	}
	if got := u.Query()["section"]; !slices.Equal(got, []string{"101", "PIT B"}) {
		t.Errorf("section %q", got)
	}
	if got := u.Query()["kinds"]; !slices.Equal(got, []string{"restock,sale"}) {
		t.Errorf("kinds %q", got)
	}
	// String path parameters are escaped.
	if _, err := c.GetSavedList(t.Context(), "a b/c"); err != nil {
		t.Fatal(err)
	}
	if got := rec.requests[1].URL.EscapedPath(); got != "/v1/lists/a%20b%2Fc" {
		t.Errorf("escaped path %s", got)
	}
}

func TestJSONBody(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{{status: 201, body: map[string]any{"id": "l1", "name": "Tour"}}})
	list, err := c.CreateSavedList(t.Context(), &CreateSavedListRequest{Name: "Tour"})
	if err != nil {
		t.Fatal(err)
	}
	if list.ID != "l1" {
		t.Errorf("list %+v", list)
	}
	req := rec.requests[0]
	if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("%s %q", req.Method, req.Header.Get("Content-Type"))
	}
	if got := string(rec.bodies[0]); got != `{"name":"Tour"}` {
		t.Errorf("body %s", got)
	}
}

func TestNullableClearsAndOmits(t *testing.T) {
	data, err := json.Marshal(WebhookPatch{ListID: Null[string]()})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"listId":null}` {
		t.Errorf("null: %s", data)
	}
	data, _ = json.Marshal(WebhookPatch{})
	if string(data) != `{}` {
		t.Errorf("unset: %s", data)
	}
	var p WebhookPatch
	if err := json.Unmarshal([]byte(`{"listId":"l1"}`), &p); err != nil {
		t.Fatal(err)
	}
	if v, ok := p.ListID.Get(); !ok || v != "l1" {
		t.Errorf("decoded %+v", p.ListID)
	}
}

func TestRetries429AfterRetryAfter(t *testing.T) {
	c, rec, sleeps := newTestClient(t, []reply{
		problem(429, map[string]string{"Retry-After": "2"}, map[string]any{"type": "t", "title": "Rate limit exceeded"}),
		emptyPage(),
	})
	if _, err := c.ListEvents(t.Context(), nil); err != nil {
		t.Fatal(err)
	}
	if len(rec.requests) != 2 || !slices.Equal(*sleeps, []time.Duration{2 * time.Second}) {
		t.Errorf("requests %d sleeps %v", len(rec.requests), *sleeps)
	}
}

func TestRetries429ForPost(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{
		problem(429, map[string]string{"Retry-After": "1"}, map[string]any{"type": "t", "title": "Rate limit exceeded"}),
		{status: 201, body: map[string]any{"id": "l1"}},
	})
	if _, err := c.CreateSavedList(t.Context(), &CreateSavedListRequest{Name: "x"}); err != nil {
		t.Fatal(err)
	}
	if len(rec.requests) != 2 {
		t.Errorf("requests %d", len(rec.requests))
	}
}

func TestRateLimitErrorWhenRetryAfterTooLong(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{problem(429, map[string]string{
		"Retry-After":      "3600",
		"RateLimit":        `"day";r=0;t=3600`,
		"RateLimit-Policy": `"minute";q=60;w=60, "day";q=50000;w=86400`,
	}, map[string]any{"type": "https://status302.com/problems/rate-limited", "title": "Rate limit exceeded", "limit": "day"})})
	_, err := c.ListEvents(t.Context(), nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !errors.Is(err, ErrRateLimited) {
		t.Fatalf("err %v", err)
	}
	if !apiErr.HasRetryAfter || apiErr.RetryAfter != time.Hour || apiErr.Limit() != "day" {
		t.Errorf("retry after %v %v limit %q", apiErr.HasRetryAfter, apiErr.RetryAfter, apiErr.Limit())
	}
	rl := apiErr.RateLimit
	if rl == nil || rl.Policy != "day" || rl.Reset != 3600 || rl.Remaining != 0 {
		t.Fatalf("rate limit %+v", rl)
	}
	if want := []RateLimitPolicy{{"minute", 60, 60}, {"day", 50000, 86400}}; !slices.Equal(rl.Policies, want) {
		t.Errorf("policies %+v", rl.Policies)
	}
	if len(rec.requests) != 1 {
		t.Errorf("requests %d", len(rec.requests))
	}
}

func TestRetries5xxAndTransportErrorsForGet(t *testing.T) {
	c, rec, sleeps := newTestClient(t, []reply{
		problem(503, nil, map[string]any{"type": "about:blank", "title": "Service Unavailable"}),
		{err: errors.New("connection reset")},
		problem(502, nil, map[string]any{"type": "about:blank", "title": "Bad Gateway"}),
		problem(500, nil, map[string]any{"type": "about:blank", "title": "Internal Server Error"}),
	})
	_, err := c.ListEvents(t.Context(), nil)
	if !errors.Is(err, ErrServer) {
		t.Fatalf("err %v", err)
	}
	if len(rec.requests) != 4 || len(*sleeps) != 3 {
		t.Fatalf("requests %d sleeps %v", len(rec.requests), *sleeps)
	}
	for i, s := range *sleeps {
		if s < 0 || s > baseDelay<<i {
			t.Errorf("sleep %d = %v beyond %v", i, s, baseDelay<<i)
		}
	}
}

func TestNoRetryOf500OnPost(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{problem(500, nil, map[string]any{"type": "about:blank", "title": "Internal Server Error"})})
	if _, err := c.CreateSavedList(t.Context(), &CreateSavedListRequest{Name: "x"}); !errors.Is(err, ErrServer) {
		t.Fatalf("err %v", err)
	}
	if len(rec.requests) != 1 {
		t.Errorf("requests %d", len(rec.requests))
	}
}

func TestConnectionErrorAfterRetries(t *testing.T) {
	cause := errors.New("dial tcp: refused")
	c, rec, _ := newTestClient(t, []reply{{err: cause}, {err: cause}, {err: cause}, {err: cause}})
	_, err := c.ListEvents(t.Context(), nil)
	var connErr *ConnectionError
	if !errors.As(err, &connErr) || !errors.Is(err, cause) || connErr.Method != http.MethodGet {
		t.Fatalf("err %v", err)
	}
	if len(rec.requests) != 4 {
		t.Errorf("requests %d", len(rec.requests))
	}
}

func TestTypedErrors(t *testing.T) {
	c, _, _ := newTestClient(t, []reply{
		problem(402, nil, map[string]any{"type": "https://status302.com/problems/plan-required", "title": "Plan required", "requiredPlan": "elite", "currentPlan": "pro", "feature": "firehose"}),
		problem(403, nil, map[string]any{"type": "https://status302.com/problems/insufficient-scope", "title": "Insufficient scope", "requiredScope": "market:read"}),
		problem(401, nil, map[string]any{"type": "https://status302.com/problems/invalid-api-key", "title": "Invalid API key"}),
		problem(422, nil, map[string]any{"type": "about:blank", "title": "Invalid parameter"}),
		problem(409, nil, map[string]any{"type": "about:blank", "title": "Conflict"}),
		{status: 418, body: nil},
	})
	var apiErr *APIError
	_, err := c.ListEvents(t.Context(), nil)
	if !errors.Is(err, ErrPlanRequired) || !errors.As(err, &apiErr) {
		t.Fatalf("402: %v", err)
	}
	if apiErr.RequiredPlan() != "elite" || apiErr.CurrentPlan() != "pro" || apiErr.Feature() != "firehose" {
		t.Errorf("402 extensions %q %q %q", apiErr.RequiredPlan(), apiErr.CurrentPlan(), apiErr.Feature())
	}
	_, err = c.ListEvents(t.Context(), nil)
	if !errors.Is(err, ErrPermission) || !errors.As(err, &apiErr) || apiErr.RequiredScope() != "market:read" || apiErr.Type() != "https://status302.com/problems/insufficient-scope" {
		t.Errorf("403: %v", err)
	}
	if _, err = c.ListEvents(t.Context(), nil); !errors.Is(err, ErrAuthentication) || errors.Is(err, ErrPermission) {
		t.Errorf("401: %v", err)
	}
	if _, err = c.ListEvents(t.Context(), nil); !errors.Is(err, ErrBadRequest) {
		t.Errorf("422: %v", err)
	}
	if _, err = c.ListEvents(t.Context(), nil); !errors.Is(err, ErrConflict) {
		t.Errorf("409: %v", err)
	}
	// A body that is not problem+json gets a problem from the status.
	_, err = c.ListEvents(t.Context(), nil)
	if !errors.As(err, &apiErr) || apiErr.Problem.Title != "I'm a teapot" || apiErr.Problem.Type != "about:blank" || apiErr.StatusCode != http.StatusTeapot {
		t.Errorf("418: %v %+v", err, apiErr)
	}
}

func TestErrorsCarryRequestID(t *testing.T) {
	c, _, _ := newTestClient(t, []reply{
		problem(404, map[string]string{"X-Request-Id": "req_abc"}, map[string]any{"type": "about:blank", "title": "Not Found", "detail": "list not found"}),
		problem(404, nil, map[string]any{"type": "about:blank", "title": "Not Found"}),
	})
	var apiErr *APIError
	_, err := c.ListEvents(t.Context(), nil)
	if !errors.Is(err, ErrNotFound) || !errors.As(err, &apiErr) || apiErr.RequestID != "req_abc" || err.Error() != "Not Found: list not found (request req_abc)" {
		t.Errorf("tagged: %v", err)
	}
	_, err = c.ListEvents(t.Context(), nil)
	if !errors.As(err, &apiErr) || apiErr.RequestID != "" || err.Error() != "Not Found" {
		t.Errorf("bare: %v", err)
	}
}

func TestPaginationFollowsNextCursor(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{
		ok(map[string]any{"items": []any{map[string]any{"id": 1}, map[string]any{"id": 2}}, "nextCursor": "c2"}, nil),
		ok(map[string]any{"items": []any{map[string]any{"id": 3}}, "nextCursor": nil}, nil),
	})
	var ids []int64
	for event, err := range c.ListEventsAll(t.Context(), &ListEventsParams{Limit: new(2)}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, event.ID)
	}
	if !slices.Equal(ids, []int64{1, 2, 3}) {
		t.Errorf("ids %v", ids)
	}
	if want := []string{"http://api.test/v1/events?limit=2", "http://api.test/v1/events?cursor=c2&limit=2"}; !slices.Equal(rec.urls(), want) {
		t.Errorf("urls %v", rec.urls())
	}
}

func TestPaginationFollowsNextBeforeForAlerts(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{
		ok(map[string]any{"items": []any{map[string]any{"id": 9}, map[string]any{"id": 8}}, "nextBefore": 8, "usage": map[string]any{"rules": 3, "limit": 10}}, nil),
		ok(map[string]any{"items": []any{map[string]any{"id": 7}}, "usage": map[string]any{"rules": 3, "limit": 10}}, nil),
		ok(map[string]any{"items": []any{map[string]any{"id": 30}}, "nextBefore": 30}, nil),
		ok(map[string]any{"items": []any{}}, nil),
	})
	var ids []int64
	for rule, err := range c.ListAlertRulesAll(t.Context(), &ListAlertRulesParams{Limit: new(2)}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, rule.ID)
	}
	if !slices.Equal(ids, []int64{9, 8, 7}) {
		t.Errorf("ids %v", ids)
	}
	var befores []*int64
	for page, err := range c.ListAlertEventsPages(t.Context(), nil) {
		if err != nil {
			t.Fatal(err)
		}
		befores = append(befores, page.NextBefore)
	}
	if len(befores) != 2 || befores[0] == nil || *befores[0] != 30 || befores[1] != nil {
		t.Errorf("nextBefore %v", befores)
	}
	want := []string{
		"http://api.test/v1/alerts/rules?limit=2",
		"http://api.test/v1/alerts/rules?before=8&limit=2",
		"http://api.test/v1/alerts/events",
		"http://api.test/v1/alerts/events?before=30",
	}
	if !slices.Equal(rec.urls(), want) {
		t.Errorf("urls %v", rec.urls())
	}
}

func TestPaginationStopsOnError(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{
		ok(map[string]any{"items": []any{map[string]any{"id": 1}}, "nextCursor": "c2"}, nil),
		problem(404, nil, map[string]any{"type": "about:blank", "title": "Not Found"}),
	})
	var ids []int64
	var errs []error
	for event, err := range c.ListEventsAll(t.Context(), nil) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		ids = append(ids, event.ID)
	}
	if !slices.Equal(ids, []int64{1}) || len(errs) != 1 || !errors.Is(errs[0], ErrNotFound) || len(rec.requests) != 2 {
		t.Errorf("ids %v errs %v requests %d", ids, errs, len(rec.requests))
	}
}

func TestIfNoneMatchAndETag(t *testing.T) {
	event := map[string]any{"id": 46, "name": "Tour"}
	c, rec, _ := newTestClient(t, []reply{
		ok(event, map[string]string{"ETag": `W/"v1"`}),
		notModified(`W/"v1"`, map[string]string{"RateLimit": `"minute";r=58;t=40`, "RateLimit-Policy": `"minute";q=60;w=60`}),
		ok(map[string]any{"id": 46, "name": "Tour II"}, map[string]string{"ETag": `W/"v2"`}),
	})
	var first Response
	detail, err := c.GetEvent(t.Context(), 46, WithResponseInto(&first))
	if err != nil {
		t.Fatal(err)
	}
	if first.StatusCode != http.StatusOK || first.ETag != `W/"v1"` || detail.Name != "Tour" || rec.requests[0].Header.Get("If-None-Match") != "" {
		t.Errorf("first %+v %+v", first, detail)
	}

	_, err = c.GetEvent(t.Context(), 46, WithIfNoneMatch(first.ETag))
	var nm *NotModifiedError
	if !errors.Is(err, ErrNotModified) || !errors.As(err, &nm) {
		t.Fatalf("not modified: %v", err)
	}
	if rec.requests[1].Header.Get("If-None-Match") != `W/"v1"` || nm.ETag != `W/"v1"` || nm.RateLimit == nil || nm.RateLimit.Remaining != 58 {
		t.Errorf("304 %+v", nm)
	}
	var var2 *APIError
	if errors.As(err, &var2) {
		t.Error("a 304 is not an APIError")
	}

	var changed Response
	detail, err = c.GetEvent(t.Context(), 46, WithIfNoneMatch(first.ETag), WithResponseInto(&changed))
	if err != nil || changed.StatusCode != http.StatusOK || changed.ETag != `W/"v2"` || detail.Name != "Tour II" {
		t.Errorf("changed %v %+v %+v", err, changed, detail)
	}
}

func TestETagCacheReturnsCachedBodyOn304(t *testing.T) {
	cache := NewMemoryETagCache(0)
	page := map[string]any{"items": []any{map[string]any{"id": 1}}, "nextCursor": nil}
	c, rec, _ := newTestClient(t, []reply{
		ok(page, map[string]string{"ETag": `W/"p1"`}),
		notModified(`W/"p1"`, nil),
		notModified(`W/"p1"`, nil),
		ok(page, map[string]string{"ETag": `W/"p2"`}),
	}, WithETagCache(cache))
	first, err := c.ListEvents(t.Context(), &ListEventsParams{Limit: new(1)})
	if err != nil || len(first.Items) != 1 || cache.Len() != 1 {
		t.Fatalf("first %v %+v len %d", err, first, cache.Len())
	}
	var hit Response
	again, err := c.ListEvents(t.Context(), &ListEventsParams{Limit: new(1)}, WithResponseInto(&hit))
	if err != nil || rec.requests[1].Header.Get("If-None-Match") != `W/"p1"` || hit.StatusCode != http.StatusNotModified || hit.ETag != `W/"p1"` {
		t.Fatalf("hit %v %+v", err, hit)
	}
	if len(again.Items) != 1 || again.Items[0].ID != 1 || again == first {
		t.Errorf("cached page %+v", again)
	}
	// A caller-chosen condition bypasses the cache and gets NotModified.
	if _, err := c.GetEvent(t.Context(), 1, WithIfNoneMatch(`W/"p1"`)); !errors.Is(err, ErrNotModified) {
		t.Errorf("explicit condition: %v", err)
	}
	// Another URL is another entry.
	if _, err := c.ListEvents(t.Context(), &ListEventsParams{Limit: new(2)}); err != nil {
		t.Fatal(err)
	}
	if rec.requests[3].Header.Get("If-None-Match") != "" || cache.Len() != 2 {
		t.Errorf("second entry: %q len %d", rec.requests[3].Header.Get("If-None-Match"), cache.Len())
	}
}

func TestMemoryETagCacheEvictsLeastRecentlyUsed(t *testing.T) {
	ctx := t.Context()
	cache := NewMemoryETagCache(2)
	cache.Set(ctx, "a", CachedResponse{ETag: "1"})
	cache.Set(ctx, "b", CachedResponse{ETag: "2"})
	cache.Get(ctx, "a")
	cache.Set(ctx, "c", CachedResponse{ETag: "3"})
	var got []string
	for _, k := range []string{"a", "b", "c"} {
		e, ok := cache.Get(ctx, k)
		if !ok {
			got = append(got, "")
			continue
		}
		got = append(got, e.ETag)
	}
	if !slices.Equal(got, []string{"1", "", "3"}) {
		t.Errorf("entries %q", got)
	}
}

func TestResponseExposesRateLimit(t *testing.T) {
	c, _, _ := newTestClient(t, []reply{ok(map[string]any{"items": []any{}, "nextCursor": nil}, map[string]string{
		"RateLimit": `"minute";r=59;t=42`, "RateLimit-Policy": `"minute";q=60;w=60`,
	})})
	var res Response
	if _, err := c.ListEvents(t.Context(), nil, WithResponseInto(&res)); err != nil {
		t.Fatal(err)
	}
	if res.RateLimit == nil || res.RateLimit.Remaining != 59 || res.RateLimit.Reset != 42 {
		t.Errorf("rate limit %+v", res.RateLimit)
	}
}

func TestRawBodiesAndNoContent(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{
		{status: 200, header: map[string]string{"Content-Type": "text/csv"}, body: "id,name\n"},
		{status: 204},
	})
	csv, err := c.ExportPortfolioCSV(t.Context(), nil)
	if err != nil || !strings.Contains(string(csv), "id,name") {
		t.Errorf("csv %q %v", csv, err)
	}
	if err := c.DeleteSavedList(t.Context(), "l1"); err != nil {
		t.Fatal(err)
	}
	if rec.requests[1].Method != http.MethodDelete {
		t.Errorf("method %s", rec.requests[1].Method)
	}
}

func TestGetPath(t *testing.T) {
	c, rec, _ := newTestClient(t, []reply{ok(map[string]any{"items": []any{}}, nil)})
	raw, err := c.GetPath(t.Context(), "/v1/events/46/metrics?scope=all")
	if err != nil || string(raw) != `{"items":[]}` {
		t.Fatalf("%s %v", raw, err)
	}
	if rec.urls()[0] != "http://api.test/v1/events/46/metrics?scope=all" {
		t.Errorf("url %s", rec.urls()[0])
	}
	if _, err := c.GetPath(t.Context(), "v1/events"); err == nil {
		t.Error("relative path accepted")
	}
}

func TestParseRetryAfter(t *testing.T) {
	cases := []struct {
		value string
		now   time.Time
		want  time.Duration
		ok    bool
	}{
		{"7", time.Now(), 7 * time.Second, true},
		{"Thu, 01 Jan 1970 00:00:10 GMT", time.Unix(0, 0), 10 * time.Second, true},
		{"", time.Now(), 0, false},
		{"soon", time.Now(), 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseRetryAfter(tc.value, tc.now)
		if got != tc.want || ok != tc.ok {
			t.Errorf("ParseRetryAfter(%q) = %v, %v; want %v, %v", tc.value, got, ok, tc.want, tc.ok)
		}
	}
}

func TestParseRateLimitAbsent(t *testing.T) {
	if ParseRateLimit(http.Header{}) != nil {
		t.Error("rate limit without headers")
	}
}
