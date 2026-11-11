package status302

import (
	"container/list"
	"context"
	"sync"
)

// CachedResponse is a response body stored by an ETagCache, keyed by request URL.
type CachedResponse struct {
	ETag        string
	Body        []byte
	ContentType string
}

// ETagCache stores the bodies of GET responses for WithETagCache. Entries
// never go stale: an ETag hashes its body, so a 304 proves the stored body is
// current. Implementations must be safe for concurrent use; a networked store
// (e.g. Redis) should treat its own failures as misses.
type ETagCache interface {
	Get(ctx context.Context, key string) (CachedResponse, bool)
	Set(ctx context.Context, key string, value CachedResponse)
}

// MemoryETagCache is an in-memory, concurrency-safe ETagCache keeping the
// most recently used responses.
type MemoryETagCache struct {
	mu         sync.Mutex
	maxEntries int
	order      *list.List // front: most recently used
	entries    map[string]*list.Element
}

type memoryEntry struct {
	key   string
	value CachedResponse
}

// NewMemoryETagCache returns an LRU cache of up to maxEntries responses
// (1000 when maxEntries < 1).
func NewMemoryETagCache(maxEntries int) *MemoryETagCache {
	if maxEntries < 1 {
		maxEntries = 1000
	}
	return &MemoryETagCache{maxEntries: maxEntries, order: list.New(), entries: map[string]*list.Element{}}
}

// Len is the number of stored responses.
func (m *MemoryETagCache) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.order.Len()
}

// Get returns the response stored under key and marks it recently used.
func (m *MemoryETagCache) Get(_ context.Context, key string) (CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	el, ok := m.entries[key]
	if !ok {
		return CachedResponse{}, false
	}
	m.order.MoveToFront(el)
	return el.Value.(*memoryEntry).value, true
}

// Set stores a response, evicting the least recently used one beyond the limit.
func (m *MemoryETagCache) Set(_ context.Context, key string, value CachedResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if el, ok := m.entries[key]; ok {
		el.Value.(*memoryEntry).value = value
		m.order.MoveToFront(el)
		return
	}
	m.entries[key] = m.order.PushFront(&memoryEntry{key, value})
	if m.order.Len() > m.maxEntries {
		oldest := m.order.Back()
		m.order.Remove(oldest)
		delete(m.entries, oldest.Value.(*memoryEntry).key)
	}
}
