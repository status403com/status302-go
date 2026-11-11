package status302

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// webhookVector is the API's webhook signing test vector, shared by every SDK.
type webhookVector struct {
	Secret    string `json:"secret"`
	Timestamp int64  `json:"timestamp"`
	Body      string `json:"body"`
	Header    string `json:"header"`
}

func loadVector(t *testing.T) (webhookVector, time.Time) {
	t.Helper()
	data, err := os.ReadFile("testdata/webhook-signature.json")
	if err != nil {
		t.Fatal(err)
	}
	var v webhookVector
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return v, time.Unix(v.Timestamp, 0)
}

func TestWebhookAcceptsGoldenSignature(t *testing.T) {
	v, at := loadVector(t)
	event, err := VerifyWebhookWith([]byte(v.Body), v.Header, []string{v.Secret}, VerifyOptions{Now: at})
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != "list.sale" || !event.Test || event.ID == "" || event.CreatedAt.IsZero() {
		t.Errorf("event %+v", event)
	}
	var data struct {
		EventID int64 `json:"eventId"`
	}
	if err := event.Decode(&data); err != nil || data.EventID != 46 {
		t.Errorf("data %+v %v", data, err)
	}
}

func TestWebhookProducesServerHeader(t *testing.T) {
	v, at := loadVector(t)
	if got := SignWebhook([]byte(v.Body), v.Secret, at); got != v.Header {
		t.Errorf("header %s, want %s", got, v.Header)
	}
}

func TestWebhookRotationSecrets(t *testing.T) {
	v, at := loadVector(t)
	if _, err := VerifyWebhookWith([]byte(v.Body), v.Header, []string{"whsec_old", v.Secret}, VerifyOptions{Now: at}); err != nil {
		t.Fatal(err)
	}
}

func TestWebhookRejects(t *testing.T) {
	v, at := loadVector(t)
	cases := []struct {
		name, body, header, secret, match string
	}{
		{"changed body", v.Body + " ", v.Header, v.Secret, "no signature matches"},
		{"wrong secret", v.Body, v.Header, "whsec_wrong", "no signature matches"},
		{"no header", v.Body, "", v.Secret, "missing"},
		{"no timestamp", v.Body, "v1=abc", v.Secret, "timestamp"},
		{"non-ASCII digit", v.Body, "t=\u00b2,v1=abc", v.Secret, "timestamp"},
		{"huge timestamp", v.Body, "t=" + strings.Repeat("9", 400) + ",v1=abc", v.Secret, "timestamp"},
		{"signed timestamp", v.Body, "t=+1793836800,v1=abc", v.Secret, "timestamp"},
		{"non-hex signature", v.Body, "t=1793836800,v1=\u00e9", v.Secret, "v1 signature"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := VerifyWebhookWith([]byte(tc.body), tc.header, []string{tc.secret}, VerifyOptions{Now: at})
			if !errors.Is(err, ErrWebhookVerification) || !strings.Contains(err.Error(), tc.match) {
				t.Errorf("err %v, want %q", err, tc.match)
			}
		})
	}
}

func TestWebhookRejectsEmptySecret(t *testing.T) {
	v, at := loadVector(t)
	// HMAC accepts a zero-length key: a delivery signed with "" must not verify.
	header := SignWebhook([]byte(v.Body), "", at)
	for _, secrets := range [][]string{{""}, nil, {"", v.Secret}} {
		_, err := VerifyWebhookWith([]byte(v.Body), header, secrets, VerifyOptions{Now: at})
		if !errors.Is(err, ErrWebhookVerification) || !strings.Contains(err.Error(), "secret") {
			t.Errorf("secrets %q: %v", secrets, err)
		}
	}
}

func TestWebhookTolerance(t *testing.T) {
	v, at := loadVector(t)
	verify := func(now time.Time, tolerance time.Duration) error {
		_, err := VerifyWebhookWith([]byte(v.Body), v.Header, []string{v.Secret}, VerifyOptions{Now: now, Tolerance: tolerance})
		return err
	}
	if err := verify(at.Add(301*time.Second), 0); err == nil || !strings.Contains(err.Error(), "old") {
		t.Errorf("301 s old: %v", err)
	}
	if err := verify(at.Add(299*time.Second), 0); err != nil {
		t.Errorf("299 s old: %v", err)
	}
	if err := verify(at.Add(-301*time.Second), 0); err == nil || !strings.Contains(err.Error(), "future") {
		t.Errorf("301 s ahead: %v", err)
	}
	if err := verify(at.Add(time.Hour), 2*time.Hour); err != nil {
		t.Errorf("custom tolerance: %v", err)
	}
}
