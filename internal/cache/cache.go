package cache

import (
	"sync"
	"time"
)

// Entry holds a cached line and its expiry.
type Entry struct {
	Line    string
	Expires time.Time
}

// Cache is a simple in-memory deduplication cache keyed by line content.
type Cache struct {
	mu      sync.Mutex
	items   map[string]Entry
	ttl     time.Duration
	maxSize int
}

// New creates a Cache with the given TTL and maximum number of entries.
// A maxSize of 0 disables the size cap.
func New(ttl time.Duration, maxSize int) *Cache {
	return &Cache{
		items:   make(map[string]Entry),
		ttl:     ttl,
		maxSize: maxSize,
	}
}

// Seen returns true if key was already recorded and has not expired.
// If the key is new or expired it records it and returns false.
func (c *Cache) Seen(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	if e, ok := c.items[key]; ok && now.Before(e.Expires) {
		return true
	}

	// Evict expired entries when we are at capacity.
	if c.maxSize > 0 && len(c.items) >= c.maxSize {
		c.evictExpired(now)
		// If still at capacity, drop oldest by simply clearing (simple policy).
		if len(c.items) >= c.maxSize {
			c.items = make(map[string]Entry)
		}
	}

	c.items[key] = Entry{Line: key, Expires: now.Add(c.ttl)}
	return false
}

// evictExpired removes all entries that have passed their TTL. Must be called
// with c.mu held.
func (c *Cache) evictExpired(now time.Time) {
	for k, e := range c.items {
		if now.After(e.Expires) {
			delete(c.items, k)
		}
	}
}

// Len returns the current number of cached entries.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
