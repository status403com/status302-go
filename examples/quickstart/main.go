// Command quickstart lists events, reads an event's market metrics, receives
// a realtime message and verifies a webhook signature.
//
//	STATUS302_API_KEY=tx_live_… go run ./examples/quickstart
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/status403com/status302-go"
)

// The sample delivery from the webhooks guide, signed at a fixed time.
const (
	sampleSecret    = "whsec_5Wq9Xr2LmT8vKz3NcB7yHd4FgJ1pQs6Ua0eRtYiOoPl"
	sampleTimestamp = 1793836800
	sampleBody      = `{"id":"0b9d7c1e-2f4a-4d6b-9e8f-1a2b3c4d5e6f","type":"list.sale","createdAt":"2026-11-05T00:00:00Z","test":true,"data":{"eventId":46,"qty":2,"priceMinor":13510,"currency":"USD"}}`
	sampleHeader    = "t=1793836800,v1=8019560260f87052f8bc68db28a88f26c757d1d2cfa129b6bbe68bb47184f79a"
)

func main() {
	ctx := context.Background()
	apiURL := os.Getenv("STATUS302_API_URL")
	if apiURL == "" {
		apiURL = status302.DefaultBaseURL
	}
	client := status302.NewClient(os.Getenv("STATUS302_API_KEY"), status302.WithBaseURL(apiURL))

	// 1. The catalog: the next upcoming events.
	page, err := client.ListEvents(ctx, &status302.ListEventsParams{Limit: new(5)})
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range page.Items {
		starts := "date TBA"
		if e.StartsAt != nil {
			starts = e.StartsAt.Format(time.RFC3339)
		}
		fmt.Printf("event %d: %s (%s)\n", e.ID, e.Name, starts)
	}

	// 2. Market metrics of the first one, per scope (primary, TM resale, …).
	if len(page.Items) > 0 {
		metrics, err := client.GetEventMetrics(ctx, page.Items[0].ID)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range metrics.Items {
			fmt.Printf("metrics %s: floor %s %s, %s tickets listed\n", m.Scope, orNull(m.Floor), orNull(m.Currency), orNull(m.Inventory))
		}
	}

	// 3. The realtime stream: wait for the next alert firing of your organization.
	stream := client.Stream(status302.StreamOptions{})
	alerts := stream.Subscribe("alerts", status302.SubscribeOptions{})
	seq, err := alerts.Ready(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("stream subscribed to alerts at seq %d; waiting for a message\n", seq)
	for msg, err := range stream.Messages(ctx) {
		if err != nil {
			log.Fatal(err)
		}
		if data, ok := msg.(*status302.DataMessage); ok {
			var alert status302.AlertEvent
			if err := data.Decode(&alert); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("stream alerts #%d: %s\n", data.Seq, alert.Title)
			break
		}
	}
	_ = stream.Close()

	// 4. Webhooks: verify a delivery's Status302-Signature over its raw body.
	event, err := status302.VerifyWebhookWith([]byte(sampleBody), sampleHeader, []string{sampleSecret},
		status302.VerifyOptions{Now: time.Unix(sampleTimestamp, 0)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("webhook verified: %s %s\n", event.Type, event.ID)
}

func orNull[T any](v *T) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprint(*v)
}
