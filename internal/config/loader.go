package config

import (
	"fmt"
	"time"
)

// TimeLayout is the expected format for since/until flag values.
const TimeLayout = time.RFC3339

// ParseTime parses a time string using TimeLayout and returns an error that
// includes the field name to aid user-facing error messages.
func ParseTime(field, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(TimeLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("config: invalid %s %q — expected RFC3339 (e.g. 2006-01-02T15:04:05Z)", field, value)
	}
	return t, nil
}

// FromFlags constructs and validates a Config from raw string flag values.
// It is the bridge between the CLI flag layer and the typed Config struct.
func FromFlags(rawURL, pattern, since, until string, pretty bool, timeout time.Duration, maxLines int) (*Config, error) {
	sinceT, err := ParseTime("since", since)
	if err != nil {
		return nil, err
	}

	untilT, err := ParseTime("until", until)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		URL:      rawURL,
		Pattern:  pattern,
		Since:    sinceT,
		Until:    untilT,
		Pretty:   pretty,
		Timeout:  timeout,
		MaxLines: maxLines,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
