package sampler

import (
	"testing"
)

func TestNewValidRate(t *testing.T) {
	s, err := New(0.5, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Errorf("expected rate 0.5, got %f", s.Rate())
	}
}

func TestNewRateOne(t *testing.T) {
	s, err := New(1.0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Error("rate=1.0 should always keep lines")
		}
	}
}

func TestNewInvalidRateZero(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Error("expected error for rate=0")
	}
}

func TestNewInvalidRateNegative(t *testing.T) {
	_, err := New(-0.1, 0)
	if err == nil {
		t.Error("expected error for negative rate")
	}
}

func TestNewInvalidRateAboveOne(t *testing.T) {
	_, err := New(1.1, 0)
	if err == nil {
		t.Error("expected error for rate > 1")
	}
}

func TestKeepApproximateRate(t *testing.T) {
	s, err := New(0.2, 99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const n = 10000
	kept := 0
	for i := 0; i < n; i++ {
		if s.Keep() {
			kept++
		}
	}
	ratio := float64(kept) / n
	// Allow ±5% tolerance around 20%
	if ratio < 0.15 || ratio > 0.25 {
		t.Errorf("expected ~20%% kept, got %.2f%%", ratio*100)
	}
}

func TestKeepConcurrentSafe(t *testing.T) {
	s, err := New(0.5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				s.Keep()
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
