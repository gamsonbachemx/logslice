// Package sampler provides probabilistic log line sampling.
// It allows reducing high-volume log streams by keeping only
// a configurable fraction of lines.
package sampler

import (
	"errors"
	"math/rand"
	"sync"
)

// Sampler decides whether a given log line should be kept
// based on a configured sampling rate.
type Sampler struct {
	rate float64
	mu   sync.Mutex
	rng  *rand.Rand
}

// New creates a Sampler that keeps approximately rate*100 percent
// of lines. rate must be in the range (0, 1]. A rate of 1.0 keeps
// all lines (no sampling). A rate of 0.1 keeps ~10% of lines.
func New(rate float64, seed int64) (*Sampler, error) {
	if rate <= 0 || rate > 1 {
		return nil, errors.New("sampler: rate must be in range (0, 1]")
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)), //nolint:gosec
	}, nil
}

// Keep returns true if the line should be included in the output.
func (s *Sampler) Keep() bool {
	if s.rate == 1.0 {
		return true
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
