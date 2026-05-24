package cache

import (
	"testing"
	"time"
)

func TestSeenNewKey(t *testing.T) {
	c := New(time.Second, 0)
	if c.Seen("hello") {
		t.Fatal("expected false for unseen key")
	}
}

func TestSeenDuplicateKey(t *testing.T) {
	c := New(time.Second, 0)
	c.Seen("hello")
	if !c.Seen("hello") {
		t.Fatal("expected true for duplicate key within TTL")
	}
}

func TestSeenExpiredKey(t *testing.T) {
	c := New(10*time.Millisecond, 0)
	c.Seen("hello")
	time.Sleep(20 * time.Millisecond)
	if c.Seen("hello") {
		t.Fatal("expected false for expired key")
	}
}

func TestLenTracksEntries(t *testing.T) {
	c := New(time.Second, 0)
	c.Seen("a")
	c.Seen("b")
	c.Seen("a") // duplicate, no new entry
	if got := c.Len(); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestMaxSizeEvictsExpired(t *testing.T) {
	c := New(10*time.Millisecond, 2)
	c.Seen("x")
	c.Seen("y")
	time.Sleep(20 * time.Millisecond)
	// Adding a third entry should trigger eviction of expired items.
	c.Seen("z")
	if c.Len() > 2 {
		t.Fatalf("expected at most 2 entries after eviction, got %d", c.Len())
	}
}

func TestMaxSizeClearsWhenFull(t *testing.T) {
	// Use a very long TTL so entries never expire naturally.
	c := New(time.Hour, 2)
	c.Seen("a")
	c.Seen("b")
	// Cache is full and nothing is expired; adding a new key should reset.
	c.Seen("c")
	if c.Len() > 2 {
		t.Fatalf("expected cache reset, got %d entries", c.Len())
	}
}

func TestDistinctKeys(t *testing.T) {
	c := New(time.Second, 0)
	keys := []string{"alpha", "beta", "gamma"}
	for _, k := range keys {
		if c.Seen(k) {
			t.Fatalf("key %q should not have been seen yet", k)
		}
	}
	if c.Len() != len(keys) {
		t.Fatalf("expected %d entries, got %d", len(keys), c.Len())
	}
}
