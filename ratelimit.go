package status302

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RateLimitPolicy is one entry of the RateLimit-Policy header
// (draft-ietf-httpapi-ratelimit-headers-09, structured fields):
//
//	RateLimit-Policy: "minute";q=60;w=60, "day";q=5000;w=86400
type RateLimitPolicy struct {
	// Name is the policy name, e.g. `minute`, `day`, `key_minute`.
	Name string
	// Quota is the number of requests allowed per window.
	Quota int64
	// Window is the window length in seconds.
	Window int64
}

// RateLimitInfo is the parsed RateLimit header: it describes the most
// constrained policy (`RateLimit: "minute";r=59;t=42`).
type RateLimitInfo struct {
	// Policy is the name of the policy the counters below belong to.
	Policy string
	// Remaining is the number of requests left in the current window.
	Remaining int64
	// Reset is the number of seconds until the window resets.
	Reset int64
	// Policies lists every policy applied to the request.
	Policies []RateLimitPolicy
}

type sfItem struct {
	name   string
	params map[string]int64
}

// parseList parses a structured-field list of string items with integer parameters.
func parseList(value string) []sfItem {
	var items []sfItem
	for _, member := range strings.Split(value, ",") {
		parts := strings.Split(member, ";")
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		if len(name) >= 2 && name[0] == '"' && name[len(name)-1] == '"' {
			name = name[1 : len(name)-1]
		}
		params := map[string]int64{}
		for _, p := range parts[1:] {
			k, v, ok := strings.Cut(strings.TrimSpace(p), "=")
			if !ok || k == "" {
				continue
			}
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				params[k] = n
			}
		}
		items = append(items, sfItem{name, params})
	}
	return items
}

// ParseRateLimit reads the RateLimit and RateLimit-Policy headers; nil when the response has none.
func ParseRateLimit(h http.Header) *RateLimitInfo {
	current := parseList(h.Get("RateLimit"))
	if len(current) == 0 {
		return nil
	}
	info := &RateLimitInfo{Policy: current[0].name, Remaining: current[0].params["r"], Reset: current[0].params["t"]}
	for _, p := range parseList(h.Get("RateLimit-Policy")) {
		info.Policies = append(info.Policies, RateLimitPolicy{Name: p.name, Quota: p.params["q"], Window: p.params["w"]})
	}
	return info
}

// ParseRetryAfter reads a Retry-After value (delta-seconds or an HTTP date,
// relative to now); false when it is absent or invalid.
func ParseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if n, err := strconv.ParseUint(value, 10, 31); err == nil {
		return time.Duration(n) * time.Second, true
	}
	at, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	d := at.Sub(now)
	if d <= 0 {
		return 0, true
	}
	// Whole seconds, rounded up.
	return (d + time.Second - 1).Truncate(time.Second), true
}
