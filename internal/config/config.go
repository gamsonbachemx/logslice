package config

import (
	"errors"
	"net/url"
	"time"
)

// Config holds all runtime configuration for a logslice session.
type Config struct {
	URL     string
	Pattern string
	Since   time.Time
	Until   time.Time
	Timeout time.Duration
	Pretty  bool
	Fields  []string
	// Tail enables continuous polling mode.
	Tail         bool
	TailInterval time.Duration
	TailRetries  int
}

// Validate checks that the Config fields are consistent and returns an
// error describing the first problem found.
func Validate(c Config) error {
	if c.URL == "" {
		return errors.New("url is required")
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return errors.New("url is invalid")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url scheme must be http or https")
	}
	if !c.Until.IsZero() && !c.Since.IsZero() && c.Until.Before(c.Since) {
		return errors.New("until must not be before since")
	}
	if c.Timeout < 0 {
		return errors.New("timeout must not be negative")
	}
	if c.Tail && c.TailInterval <= 0 {
		return errors.New("tail-interval must be positive when tail is enabled")
	}
	if c.Tail && c.TailRetries <= 0 {
		return errors.New("tail-retries must be positive when tail is enabled")
	}
	return nil
}
