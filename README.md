# status302-go — Go SDK

Official Go client for the [Status302 API](https://status302.com/docs/): real-time ticket market data (catalog, metrics, sales, seat maps) over REST, a realtime WebSocket stream and signed webhooks.

- One typed method per public operation, with structs and enums generated from the public OpenAPI spec
- Automatic retries on `429` (honouring `Retry-After`) and on `5xx`/connection errors for idempotent methods; parsed `RateLimit` headers
- Pagination iterators (`range` over `ListEventsAll`, `ListEventsPages`, …; cursors and alert `nextBefore` ids) and typed errors (`errors.Is(err, status302.ErrPlanRequired)`, `*APIError` with the problem+json body and the response's request id, also appended to the message as `(request <id>)`; quote it to support)
- Conditional requests: `WithIfNoneMatch` with `ErrNotModified`, response ETags, and an opt-in ETag cache
- Stream client: bearer or ticket auth, reconnects with backoff, `since` resume per topic, resync notices with a snapshot helper, seq deduplication, JSON or MessagePack frames
- Webhook signature verification

The module version equals the version of the public API spec it was generated from. Full guide: <https://status302.com/docs/sdks/go/>.

## Install

```sh
go get github.com/status403com/status302-go
```

Go ≥ 1.26. Dependencies: [coder/websocket](https://github.com/coder/websocket) and [vmihailenco/msgpack](https://github.com/vmihailenco/msgpack) (stream client only).

## Quickstart

```go
import "github.com/status403com/status302-go"

client := status302.NewClient(os.Getenv("STATUS302_API_KEY"))

page, err := client.ListEvents(ctx, &status302.ListEventsParams{Country: new("US"), Limit: new(5)})
if err != nil {
	log.Fatal(err)
}
for _, e := range page.Items {
	fmt.Println(e.ID, e.Name, e.StartsAt)
}

metrics, err := client.GetEventMetrics(ctx, page.Items[0].ID)
if err != nil {
	log.Fatal(err)
}
for _, m := range metrics.Items {
	fmt.Println(m.Scope, m.Floor, m.Currency) // money is integer minor units
}
```

Create an API key in the app under **Settings → Developers**. Keys look like `tx_live_<keyId>_<secret>` and are sent as `Authorization: Bearer …`.

Methods are named after the operation ids (`ListEvents`, `GetEventMetrics`, `CreateSavedList`, …). Path parameters are positional, query parameters go in a `<Operation>Params` struct (nil for none), and request bodies are pointers to the generated structs. Optional fields are pointers (`new(value)`); optional fields that may also be null are `Nullable[T]` (`status302.Value(v)`, `status302.Null[T]()`, unset = omitted), so a PATCH can clear a property. CSV and binary downloads return `[]byte`.

Client options: `WithBaseURL` (default `https://api.status302.com`), `WithMaxRetries` (3), `WithMaxRetryAfter` (60 s), `WithTimeout` (30 s per attempt), `WithHTTPClient` and `WithETagCache`. Per-call options: `WithHeader`, `WithIfNoneMatch` and `WithResponseInto(&resp)` for the status, headers, ETag and parsed rate limit.

### Pagination

```go
for event, err := range client.ListEventsAll(ctx, &status302.ListEventsParams{Country: new("US"), Limit: new(200)}) {
	if err != nil {
		return err
	}
	// …
}
```

`…Pages` yields whole pages instead (save `NextCursor` to resume). Both follow `nextCursor` until it is null, or `nextBefore` (as `before`) for `ListAlertRules` and `ListAlertEvents`.

### Conditional requests

```go
var first status302.Response
metrics, err := client.GetEventMetrics(ctx, 4102, status302.WithResponseInto(&first)) // first.ETag

_, err = client.GetEventMetrics(ctx, 4102, status302.WithIfNoneMatch(first.ETag))
if errors.Is(err, status302.ErrNotModified) {
	fmt.Println("unchanged")
}

cached := status302.NewClient(key, status302.WithETagCache(status302.NewMemoryETagCache(1000)))
```

Operations that accept `If-None-Match` say so in their doc comment; a `304` returns a `*NotModifiedError` (with the ETag and rate limit state). With `WithETagCache` the client revalidates every `GET` by itself and returns the stored body on `304`; `MemoryETagCache` is LRU and concurrency-safe, and any `ETagCache` (`Get`/`Set`) works, e.g. one backed by Redis.

### Errors and rate limits

```go
_, err := client.ListMovers(ctx, &status302.ListMoversParams{Window: new(status302.ListMoversWindow24h)})
var apiErr *status302.APIError
switch {
case errors.Is(err, status302.ErrPlanRequired) && errors.As(err, &apiErr):
	fmt.Println("needs plan", apiErr.RequiredPlan())
case errors.Is(err, status302.ErrRateLimited) && errors.As(err, &apiErr):
	fmt.Println("retry in", apiErr.RetryAfter, "quota:", apiErr.RateLimit)
case err != nil:
	return err
}
```

| Sentinel | Status |
|---|---|
| `ErrBadRequest` | `400`, `422` |
| `ErrAuthentication` | `401` |
| `ErrPlanRequired` | `402` (`RequiredPlan()`, `CurrentPlan()`, `Feature()`, `Limit()`) |
| `ErrPermission` | `403` (`RequiredScope()`) |
| `ErrNotFound` | `404` |
| `ErrConflict` | `409` |
| `ErrRateLimited` | `429` (`RetryAfter`, `RateLimit`, `Limit()`) |
| `ErrServer` | `5xx` |

Network failures return a `*ConnectionError` wrapping the cause.

### Realtime stream

```go
stream := client.Stream(status302.StreamOptions{})
defer stream.Close()

alerts := stream.Subscribe("alerts", status302.SubscribeOptions{})
seq, err := alerts.Ready(ctx) // head seq, or a *SubscriptionError
for msg, err := range stream.Messages(ctx) {
	if err != nil {
		return err // *StreamClosedError: 4401/4403 are final
	}
	switch m := msg.(type) {
	case *status302.DataMessage:
		var alert status302.AlertEvent
		_ = m.Decode(&alert)
		fmt.Println(m.Seq, alert.Title)
	case *status302.ResyncMessage:
		// The gap could not be replayed: reload the REST snapshot, then continue live.
		snapshot, err := stream.FetchSnapshot(ctx, m)
		_, _ = snapshot, err
	case *status302.ErrorMessage:
		fmt.Println("refused", m.Topic, m.Code)
	}
}
```

The key needs the `stream` scope to connect, and each topic the scope of its REST equivalent (`event.{id}.kpis` → `market:read`, `sales`/`event.{id}.sales` → `sales:read`, `event.{id}.primary` → `maps:read` or `market:read`, `onsales` → `catalog:read`, `alerts` → `alerts:write`); without it the subscription fails with a `forbidden_topic` `ErrorMessage`. The key is sent on the upgrade request (`StreamAuthBearer`, default with a key); `Auth: status302.StreamAuthTicket` with `GetTicket` connects with single-use tickets instead (same scopes as the key).

#### Binary encoding (MessagePack)

`StreamOptions{Encoding: status302.StreamEncodingMsgpack}` asks the server for binary MessagePack frames (`encoding=msgpack` on the URL) instead of JSON text. Messages and `Data` payloads decode to the same values as with JSON; the client's own messages stay JSON. Default: JSON.

### Webhooks

```go
body, _ := io.ReadAll(r.Body)
event, err := status302.VerifyWebhook(body, r.Header.Get("Status302-Signature"), os.Getenv("STATUS302_WEBHOOK_SECRET"))
if err != nil { // errors.Is(err, status302.ErrWebhookVerification)
	http.Error(w, "invalid signature", http.StatusBadRequest)
	return
}
```

Verify over the raw body before parsing it; pass several secrets while rotating. `SignWebhook(body, secret, time.Now())` signs test deliveries.

## Versioning

Releases share the major and minor version of the API spec they were generated from (`status302.Version`); `/v1` changes are additive only, so new fields, operations and enum values arrive in minor releases. See the [changelog](https://status302.com/docs/changelog/).

## License

[MIT](LICENSE)
