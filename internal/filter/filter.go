package filter

import (
	"encoding/json"
	"regexp"
	"time"
)

// LogEntry represents a single structured JSON log line.
type LogEntry struct {
	Timestamp time.Time
	Raw       map[string]interface{}
	Line      string
}

// Options holds the filtering criteria.
type Options struct {
	Pattern   *regexp.Regexp
	TimeStart *time.Time
	TimeEnd   *time.Time
	TimeField string
}

// ParseLine parses a raw JSON log line into a LogEntry.
func ParseLine(line string, timeField string) (*LogEntry, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, err
	}

	entry := &LogEntry{Raw: raw, Line: line}

	if tf := timeField; tf != "" {
		if v, ok := raw[tf]; ok {
			if s, ok := v.(string); ok {
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					entry.Timestamp = t
				}
			}
		}
	}

	return entry, nil
}

// Match returns true if the log entry satisfies all filter options.
func Match(entry *LogEntry, opts Options) bool {
	if opts.Pattern != nil && !opts.Pattern.MatchString(entry.Line) {
		return false
	}

	if !entry.Timestamp.IsZero() {
		if opts.TimeStart != nil && entry.Timestamp.Before(*opts.TimeStart) {
			return false
		}
		if opts.TimeEnd != nil && entry.Timestamp.After(*opts.TimeEnd) {
			return false
		}
	}

	return true
}
