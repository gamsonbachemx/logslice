package fetch

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// RetryConfig holds parameters for retry behaviour on transient HTTP errors.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Backoff is the initial wait duration, doubled on each subsequent attempt.
	Backoff time.Duration
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Backoff:     500 * time.Millisecond,
	}
}

// isRetryable reports whether the given HTTP status code warrants a retry.
func isRetryable(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

// withRetry executes fn up to cfg.MaxAttempts times, retrying on transient
// errors. It respects context cancellation between attempts.
func withRetry(ctx context.Context, cfg RetryConfig, fn func() (int, error)) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	backoff := cfg.Backoff
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		statusCode, err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Only retry on retryable HTTP status codes; propagate other errors.
		if statusCode != 0 && !isRetryable(statusCode) {
			return lastErr
		}

		// Last attempt — do not sleep.
		if attempt == cfg.MaxAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff *= 2
	}

	return errors.New("max retry attempts reached: " + lastErr.Error())
}
