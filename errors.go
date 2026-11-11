package status302

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Sentinels matched by *APIError with errors.Is, by HTTP status:
//
//	if errors.Is(err, status302.ErrPlanRequired) { … }
//
// Use errors.As with *APIError for the problem details and response headers.
var (
	// ErrBadRequest: 400 or 422, the request is malformed or a parameter is invalid.
	ErrBadRequest = errors.New("status302: bad request")
	// ErrAuthentication: 401, missing, malformed, expired or revoked credentials.
	ErrAuthentication = errors.New("status302: authentication failed")
	// ErrPlanRequired: 402, the organization's plan does not include the request.
	ErrPlanRequired = errors.New("status302: plan required")
	// ErrPermission: 403, the key lacks a scope, the client IP is not allowed, or the role is insufficient.
	ErrPermission = errors.New("status302: permission denied")
	// ErrNotFound: 404, the resource does not exist or is not visible to the caller.
	ErrNotFound = errors.New("status302: not found")
	// ErrConflict: 409, the request conflicts with the resource's current state.
	ErrConflict = errors.New("status302: conflict")
	// ErrRateLimited: 429, a quota is exhausted (see APIError.RetryAfter).
	ErrRateLimited = errors.New("status302: rate limited")
	// ErrServer: 5xx, the API failed; retried automatically for idempotent requests.
	ErrServer = errors.New("status302: server error")
	// ErrNotModified is matched by *NotModifiedError (a 304 answer to WithIfNoneMatch).
	ErrNotModified = errors.New("status302: not modified")
)

// APIError is a non-2xx answer of the API. Problem holds the problem+json
// body; the message ends with `(request <id>)` when the response names its
// request: quote it to support, it matches the server logs.
type APIError struct {
	StatusCode int
	// Problem is the RFC 9457 body (synthesized from the status when the body is not problem+json).
	Problem Problem
	// Body is the raw response body; the extension members (requiredPlan, …) are read from it.
	Body   []byte
	Header http.Header
	// RateLimit is the rate limit state from the response headers (nil when absent).
	RateLimit *RateLimitInfo
	// RequestID is the response's X-Request-Id ("" when absent).
	RequestID string
	// RetryAfter is the Retry-After header; HasRetryAfter reports whether it was present.
	RetryAfter    time.Duration
	HasRetryAfter bool
}

func newAPIError(res *httpReply, body []byte, rateLimit *RateLimitInfo, retryAfter time.Duration, hasRetryAfter bool) *APIError {
	e := &APIError{
		StatusCode: res.StatusCode, Body: body, Header: res.Header, RateLimit: rateLimit,
		RequestID: res.Header.Get("X-Request-Id"), RetryAfter: retryAfter, HasRetryAfter: hasRetryAfter,
	}
	if err := json.Unmarshal(body, &e.Problem); err != nil || e.Problem.Title == "" {
		title := http.StatusText(res.StatusCode)
		if title == "" {
			title = fmt.Sprintf("HTTP %d", res.StatusCode)
		}
		e.Problem = Problem{Type: "about:blank", Title: title, Status: res.StatusCode}
		if len(body) > 0 {
			detail := string(body)
			if len(detail) > 500 {
				detail = detail[:500]
			}
			e.Problem.Detail = &detail
		}
		e.Body = nil
	}
	if e.Problem.Type == "" {
		e.Problem.Type = "about:blank"
	}
	if e.Problem.Status == 0 {
		e.Problem.Status = res.StatusCode
	}
	return e
}

func (e *APIError) Error() string {
	msg := e.Problem.Title
	if e.Problem.Detail != nil && *e.Problem.Detail != "" {
		msg += ": " + *e.Problem.Detail
	}
	if e.RequestID != "" {
		msg += " (request " + e.RequestID + ")"
	}
	return msg
}

// Is matches the sentinel of the error's status (ErrNotFound for 404, …).
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrBadRequest:
		return e.StatusCode == http.StatusBadRequest || e.StatusCode == http.StatusUnprocessableEntity
	case ErrAuthentication:
		return e.StatusCode == http.StatusUnauthorized
	case ErrPlanRequired:
		return e.StatusCode == http.StatusPaymentRequired
	case ErrPermission:
		return e.StatusCode == http.StatusForbidden
	case ErrNotFound:
		return e.StatusCode == http.StatusNotFound
	case ErrConflict:
		return e.StatusCode == http.StatusConflict
	case ErrRateLimited:
		return e.StatusCode == http.StatusTooManyRequests
	case ErrServer:
		return e.StatusCode >= 500
	}
	return false
}

// Type is the problem type URI, e.g. `https://status302.com/problems/rate-limited`.
func (e *APIError) Type() string { return e.Problem.Type }

// extension reads a string extension member of the problem body.
func (e *APIError) extension(name string) string {
	var members map[string]json.RawMessage
	if json.Unmarshal(e.Body, &members) != nil {
		return ""
	}
	var s string
	if json.Unmarshal(members[name], &s) != nil {
		return ""
	}
	return s
}

// RequiredPlan is the plan that allows the request (402).
func (e *APIError) RequiredPlan() string { return e.extension("requiredPlan") }

// CurrentPlan is the caller's current plan (402).
func (e *APIError) CurrentPlan() string { return e.extension("currentPlan") }

// Feature is the gated capability, e.g. `market_data` (402).
func (e *APIError) Feature() string { return e.extension("feature") }

// RequiredScope is the scope the API key lacks (403 insufficient-scope).
func (e *APIError) RequiredScope() string { return e.extension("requiredScope") }

// Limit is the exhausted quota of a 429 (`minute`, `day`, `concurrent`,
// `stream_connections`, …) or, as a decimal string, the exhausted count limit of a 402.
func (e *APIError) Limit() string {
	var members map[string]json.RawMessage
	if json.Unmarshal(e.Body, &members) != nil || members["limit"] == nil {
		return ""
	}
	var s string
	if json.Unmarshal(members["limit"], &s) == nil {
		return s
	}
	return string(members["limit"])
}

// ConnectionError means the request never produced an HTTP response (DNS, TLS, reset, timeout).
type ConnectionError struct {
	Method string
	URL    string
	Err    error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("status302: %s %s: %v", e.Method, e.URL, e.Err)
}

func (e *ConnectionError) Unwrap() error { return e.Err }

// NotModifiedError is the answer to a conditional request (WithIfNoneMatch)
// while the copy with that ETag is still current (HTTP 304).
type NotModifiedError struct {
	// ETag is the current ETag, equal to (one of) the ETags sent.
	ETag      string
	Header    http.Header
	RateLimit *RateLimitInfo
}

func (e *NotModifiedError) Error() string { return "status302: not modified (ETag " + e.ETag + ")" }

// Is matches ErrNotModified.
func (e *NotModifiedError) Is(target error) bool { return target == ErrNotModified }
