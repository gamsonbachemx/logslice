// Package cache provides a lightweight, thread-safe in-memory deduplication
// cache for log lines streamed by logslice.
//
// Entries are keyed by raw line content and expire after a configurable TTL.
// An optional maximum size prevents unbounded memory growth by evicting expired
// entries or, as a last resort, resetting the cache when it is full and no
// expired entries exist.
//
// Typical usage:
//
//	c := cache.New(30*time.Second, 10_000)
//	if !c.Seen(line) {
//		// process line
//	}
package cache
