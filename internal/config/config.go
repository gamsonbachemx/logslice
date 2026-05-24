// Package config defines the runtime configuration for logslice and provides
// validation helpers.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Config holds all runtime settings derived from CLI flags or other sources.
type Config struct {
	URL     string
	Pattern string
	Since   time.Time
	Until   time.Time
	Timeout time.Duration
	Pretty  bool
	Rate    float64
	// Fields is an optional comma-separated field projection spec, e.g.
	// "level,msg" or "level:severity,msg:message".
	Fields string
}

// Validate checks that the Config is internally consistent and ready for use.
func Validate(c *Config) error {
	if c.URL == "" {
		return errors.New("config: URL is required")
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("config: invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("config: URL scheme must be http or https, got %q", u.Scheme)
	}
	if !c.Since.IsZero() && !c.Until.IsZero() && c.Until.Before(c.Since) {
		return errors.New("config: --until must be after --since")
	}
	if c.Timeout < 0 {
		return errors.New("config: timeout must be non-negative")
	}
	if c.Rate < 0 {
		return errors.New("config: rate must be non-negative")
	}
	return nil
}
