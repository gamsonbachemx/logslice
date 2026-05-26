package jitter_test

import (
	"testing"
	"time"

	"github.com/yourusername/logslice/internal/jitter"
)

func TestFullZeroBase(t *testing.T) {
	if got := jitter.Full(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFullNegativeBase(t *testing.T) {
	if got := jitter.Full(-time.Second); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestWithSourceMin(t *testing.T) {
	// src always returns 0 → result should equal base exactly
	got := jitter.WithSource(time.Second, func() float64 { return 0 })
	if got != time.Second {
		t.Fatalf("expected 1s, got %v", got)
	}
}

func TestWithSourceMax(t *testing.T) {
	// src always returns 1 → result should equal base*2
	got := jitter.WithSource(time.Second, func() float64 { return 1 })
	if got != 2*time.Second {
		t.Fatalf("expected 2s, got %v", got)
	}
}

func TestEqualWithSourceMin(t *testing.T) {
	// src returns 0 → result = base/2
	got := jitter.EqualWithSource(time.Second, func() float64 { return 0 })
	if got != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", got)
	}
}

func TestEqualWithSourceMax(t *testing.T) {
	// src returns 1 → result = base*1.5 (half + half*1)
	got := jitter.EqualWithSource(time.Second, func() float64 { return 1 })
	if got != 1500*time.Millisecond {
		t.Fatalf("expected 1500ms, got %v", got)
	}
}

func TestCappedRespectsMax(t *testing.T) {
	// src always returns 1 → Full would give 2s, but cap is 1.5s
	got := jitter.Capped(time.Second, 1500*time.Millisecond)
	if got > 1500*time.Millisecond {
		t.Fatalf("value %v exceeds cap", got)
	}
}

func TestFullInRange(t *testing.T) {
	base := 100 * time.Millisecond
	for i := 0; i < 500; i++ {
		v := jitter.Full(base)
		if v < base || v >= 2*base {
			t.Fatalf("Full(%v) = %v out of [%v, %v)", base, v, base, 2*base)
		}
	}
}

func TestEqualInRange(t *testing.T) {
	base := 200 * time.Millisecond
	lo := base / 2
	hi := base + base/2
	for i := 0; i < 500; i++ {
		v := jitter.Equal(base)
		if v < lo || v > hi {
			t.Fatalf("Equal(%v) = %v out of [%v, %v]", base, v, lo, hi)
		}
	}
}
