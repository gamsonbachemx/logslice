// Package ratelimit provides a simple token-bucket rate limiter
// for controlling how many log lines are processed per second.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which log lines are consumed.
type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	done   chan struct{}
}

// New creates a Limiter that allows up to ratePerSec events per second.
// If ratePerSec is 0 or negative, no limiting is applied.
func New(ratePerSec int) *Limiter {
	l := &Limiter{
		done: make(chan struct{}),
	}
	if ratePerSec <= 0 {
		return l
	}
	interval := time.Second / time.Duration(ratePerSec)
	l.ticker = time.NewTicker(interval)
	l.tokens = make(chan struct{}, ratePerSec)
	go l.produce()
	return l
}

func (l *Limiter) produce() {
	if l.ticker == nil {
		return
	}
	for {
		select {
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
			}
		case <-l.done:
			return
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
// Returns ctx.Err() if the context is done, nil otherwise.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.tokens == nil {
		return nil
	}
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Limiter.
func (l *Limiter) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
	close(l.done)
}
