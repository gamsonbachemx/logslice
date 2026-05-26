// Package jitter provides utilities for adding randomised delay to retry
// and polling loops, preventing thundering-herd problems when many clients
// reconnect to the same upstream simultaneously.
package jitter

import (
	"math/rand"
	"time"
)

// Source is a function that returns a random float64 in [0,1).
// It is exposed so tests can inject a deterministic source.
type Source func() float64

// Jitter adds a random fraction of base to base, returning a duration in the
// range [base, base*2). If base is zero or negative the call returns zero
// immediately.
//
//	jittered := jitter.Full(time.Second)
//	time.Sleep(jittered)
func Full(base time.Duration) time.Duration {
	return WithSource(base, rand.Float64)
}

// Equal adds a random fraction of base/2 to base, returning a duration in the
// range [base*0.5, base*1.5). This keeps the centre of mass at base while
// still spreading load.
func Equal(base time.Duration) time.Duration {
	return EqualWithSource(base, rand.Float64)
}

// WithSource is like Full but accepts a custom random source, useful in tests.
func WithSource(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	return base + time.Duration(float64(base)*src())
}

// EqualWithSource is like Equal but accepts a custom random source.
func EqualWithSource(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	half := float64(base) / 2
	return time.Duration(half + half*src())
}

// Capped returns a Full jitter value that is capped at max.
func Capped(base, max time.Duration) time.Duration {
	v := Full(base)
	if v > max {
		return max
	}
	return v
}
