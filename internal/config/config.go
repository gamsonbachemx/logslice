// Package config provides configuration loading and validation for logslice.
package config

import (
	"errors"
	"net/url"
	"time"
)

// Config holds all runtime configuration for logslice.
type Config struct {
	// URL is the remote endpoint to stream logs from.
	URL string

	// Pattern is an optional regex pattern to filter log lines.
	Pattern string

	// Since filters log entries to those at or after this time.
	Since time.Time

	// Until filters log entries to those at or before this time.
	Until time.Time

	// Pretty enables human-readable output instead of raw JSON.
	Pretty bool

	// Timeout is the HTTP client timeout for the remote connection.
	Timeout time.Duration

	// MaxLines limits the number of matched lines printed (0 = unlimited).
	MaxLines int
}

// DefaultTimeout is used when no timeout is explicitly configured.
const DefaultTimeout = 30 * time.Second

// Validate checks that the Config contains valid, usable values.
func (c *Config) Validate() error {
	if c.URL == "" {
		return errors.New("config: URL must not be empty")
	}

	u, err := url.ParseRequestURI(c.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("config: URL must be a valid http or https URL")
	}

	if !c.Since.IsZero() && !c.Until.IsZero() && c.Until.Before(c.Since) {
		return errors.New("config: until must not be before since")
	}

	if c.Timeout < 0 {
		return errors.New("config: timeout must not be negative")
	}

	if c.MaxLines < 0 {
		return errors.New("config: max-lines must not be negative")
	}

	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}

	return nil
}
