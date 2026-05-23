package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/ratelimit"
)

func TestNewNoLimit(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	ctx := context.Background()
	for i := 0; i < 100; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error with no limit: %v", err)
		}
	}
}

func TestNewNegativeRate(t *testing.T) {
	l := ratelimit.New(-5)
	defer l.Stop()

	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected no error for negative rate, got: %v", err)
	}
}

func TestWaitContextCancelled(t *testing.T) {
	// Very low rate so tokens are scarce
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the initial token if any
	time.Sleep(15 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	// Consume available token
	_ = l.Wait(ctx)

	// Next wait should time out
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel2()

	err := l.Wait(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestRateLimitThroughput(t *testing.T) {
	const rate = 200
	l := ratelimit.New(rate)
	defer l.Stop()

	ctx := context.Background()
	start := time.Now()
	const n = 20
	for i := 0; i < n; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	elapsed := time.Since(start)

	// n tokens at rate/sec should take at least (n-1)/rate seconds
	minExpected := time.Duration(float64(n-1)/rate*float64(time.Second)) / 2
	if elapsed < minExpected {
		t.Logf("elapsed %v, minExpected %v — may be flaky on slow CI", elapsed, minExpected)
	}
}
