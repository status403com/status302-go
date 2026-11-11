package status302

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// WebhookEvent is a webhook delivery body. ID identifies the event and stays
// the same across retries and redeliveries (process each ID once); the
// `Status302-Delivery-Id` header identifies one delivery. Test marks
// "Send test event" deliveries.
type WebhookEvent struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	CreatedAt time.Time       `json:"createdAt"`
	Test      bool            `json:"test"`
	Data      json.RawMessage `json:"data"`
}

// Decode unmarshals the event's data into v.
func (e *WebhookEvent) Decode(v any) error { return json.Unmarshal(e.Data, v) }

// ErrWebhookVerification is matched by every error of VerifyWebhook: the
// delivery is not authentic or too old; answer it with a 4xx and do not process it.
var ErrWebhookVerification = errors.New("status302: webhook verification failed")

type webhookError struct{ msg string }

func (e *webhookError) Error() string        { return "status302: " + e.msg }
func (e *webhookError) Is(target error) bool { return target == ErrWebhookVerification }

func verifyErr(format string, args ...any) error {
	return &webhookError{fmt.Sprintf(format, args...)}
}

// VerifyOptions tunes VerifyWebhookWith.
type VerifyOptions struct {
	// Tolerance is the oldest accepted signature age (and the furthest in the future). Default 5 minutes.
	Tolerance time.Duration
	// Now is the current time (tests). Default time.Now().
	Now time.Time
}

// VerifyWebhook verifies a webhook delivery and returns its parsed body.
//
// `Status302-Signature: t=<unix seconds>,v1=<hex HMAC-SHA256(secret, "<t>.<raw body>")>`.
// Pass the raw request body exactly as received (before any JSON parsing)
// and the endpoint's signing secret (`whsec_…`); during a secret rotation
// pass both secrets. An empty secret (an unset environment variable) is
// rejected, never used as an HMAC key. Comparison is constant-time.
func VerifyWebhook(rawBody []byte, signatureHeader string, secrets ...string) (*WebhookEvent, error) {
	return VerifyWebhookWith(rawBody, signatureHeader, secrets, VerifyOptions{})
}

// VerifyWebhookWith is VerifyWebhook with a custom tolerance or clock.
func VerifyWebhookWith(rawBody []byte, signatureHeader string, secrets []string, opts VerifyOptions) (*WebhookEvent, error) {
	// HMAC accepts a zero-length key, so an empty secret would let anyone sign.
	if len(secrets) == 0 {
		return nil, verifyErr("signing secret is empty or not configured")
	}
	for _, s := range secrets {
		if s == "" {
			return nil, verifyErr("signing secret is empty or not configured")
		}
	}
	if signatureHeader == "" {
		return nil, verifyErr("missing Status302-Signature header")
	}
	var timestamp string
	var signatures [][]byte
	for _, part := range strings.Split(signatureHeader, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "t":
			timestamp = strings.TrimSpace(value)
		case "v1":
			if sig, err := hex.DecodeString(strings.TrimSpace(value)); err == nil {
				signatures = append(signatures, sig)
			}
		}
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || ts < 0 || strings.ContainsAny(timestamp, "+-") {
		return nil, verifyErr("Status302-Signature has no valid timestamp")
	}
	if len(signatures) == 0 {
		return nil, verifyErr("Status302-Signature has no v1 signature")
	}

	tolerance := opts.Tolerance
	if tolerance <= 0 {
		tolerance = 5 * time.Minute
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	age := float64(now.UnixMilli())/1000 - float64(ts)
	if age > tolerance.Seconds() {
		return nil, verifyErr("signature timestamp is %d s old (tolerance %d s)", int64(math.Round(age)), int64(tolerance.Seconds()))
	}
	if age < -tolerance.Seconds() {
		return nil, verifyErr("signature timestamp is in the future")
	}

	for _, secret := range secrets {
		expected := signature(secret, timestamp, rawBody)
		for _, sig := range signatures {
			if hmac.Equal(sig, expected) {
				var event WebhookEvent
				if err := json.Unmarshal(rawBody, &event); err != nil {
					return nil, verifyErr("body is not valid JSON: %v", err)
				}
				return &event, nil
			}
		}
	}
	return nil, verifyErr("no signature matches the signing secret")
}

// SignWebhook computes the `Status302-Signature` header for a body at
// timestamp (tests and local receivers).
func SignWebhook(rawBody []byte, secret string, timestamp time.Time) string {
	t := strconv.FormatInt(timestamp.Unix(), 10)
	return "t=" + t + ",v1=" + hex.EncodeToString(signature(secret, t, rawBody))
}

func signature(secret, timestamp string, body []byte) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return mac.Sum(nil)
}
