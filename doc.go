// Package status302 is the official Go client for the Status302 API
// (https://status302.com/docs/): real-time ticket market data over REST, a
// realtime WebSocket stream and signed webhooks.
//
// Every public operation is a Client method generated from the public
// OpenAPI spec (operations.go, types.go): path parameters are positional,
// query parameters go in a <Operation>Params struct and request bodies are
// pointers to the generated structs. Optional fields are pointers, set with
// new(value); optional fields that may also be null are Nullable.
//
//	client := status302.NewClient(os.Getenv("STATUS302_API_KEY"))
//	for event, err := range client.ListEventsAll(ctx, &status302.ListEventsParams{Country: new("US")}) {
//		if err != nil {
//			return err
//		}
//		fmt.Println(event.ID, event.Name)
//	}
//
// Errors of the API are *APIError values matching the sentinels ErrNotFound,
// ErrPlanRequired, ErrRateLimited, … with errors.Is. The client retries 429
// answers (after Retry-After) and, for idempotent methods, 5xx answers and
// connection errors.
//
// Client.Stream opens the realtime stream; VerifyWebhook checks webhook
// deliveries. Version equals the public API spec version the package was
// generated from.
package status302
