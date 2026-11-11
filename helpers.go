package status302

import (
	"bytes"
	"encoding/json"
	"iter"
	"strings"
)

// Nullable is an optional property that may also be null: unset (omitted
// from JSON), null, or a value. PATCH bodies use it to tell "leave as is"
// (unset) from "clear" (null).
type Nullable[T any] struct {
	value T
	set   bool
	null  bool
}

// Value returns a Nullable holding v.
func Value[T any](v T) Nullable[T] { return Nullable[T]{value: v, set: true} }

// Null returns a Nullable holding JSON null.
func Null[T any]() Nullable[T] { return Nullable[T]{set: true, null: true} }

// Get returns the value and whether there is one (false when unset or null).
func (n Nullable[T]) Get() (T, bool) { return n.value, n.set && !n.null }

// IsNull reports whether the property is explicitly null.
func (n Nullable[T]) IsNull() bool { return n.set && n.null }

// IsSet reports whether the property is present (a value or null).
func (n Nullable[T]) IsSet() bool { return n.set }

// IsZero reports whether the property is unset; the `omitzero` tag omits it then.
func (n Nullable[T]) IsZero() bool { return !n.set }

// MarshalJSON encodes null or the value.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.set || n.null {
		return []byte("null"), nil
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON decodes null or a value; a property absent from the JSON stays unset.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*n = Null[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*n = Value(v)
	return nil
}

// joinQuery encodes an array query parameter as comma-separated values.
func joinQuery[T any](vs []T, format func(T) string) string {
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = format(v)
	}
	return strings.Join(parts, ",")
}

// pages yields fetched pages until advance (which sets the next cursor) reports the last one.
func pages[P any](fetch func() (*P, error), advance func(*P) bool) iter.Seq2[*P, error] {
	return func(yield func(*P, error) bool) {
		for {
			page, err := fetch()
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(page, nil) || !advance(page) {
				return
			}
		}
	}
}

// items flattens pages into their items; an error ends the sequence after it is yielded.
func items[P, I any](seq iter.Seq2[*P, error], of func(*P) []I) iter.Seq2[I, error] {
	return func(yield func(I, error) bool) {
		for page, err := range seq {
			if err != nil {
				var zero I
				yield(zero, err)
				return
			}
			for _, it := range of(page) {
				if !yield(it, nil) {
					return
				}
			}
		}
	}
}
