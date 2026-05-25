// Package truncate provides field-level value truncation for structured log lines.
// It allows capping string field values to a maximum byte length, appending an
// ellipsis marker when truncation occurs.
package truncate

import (
	"encoding/json"
	"fmt"
)

const defaultEllipsis = "..."

// Truncator truncates string values in JSON log fields that exceed a maximum length.
type Truncator struct {
	fields    map[string]int // field name -> max bytes
	ellipsis  string
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithEllipsis overrides the default ellipsis string appended after truncation.
func WithEllipsis(e string) Option {
	return func(t *Truncator) {
		t.ellipsis = e
	}
}

// New creates a Truncator that caps the given fields at the specified byte lengths.
// fields is a map of JSON field name to maximum allowed byte length (excluding ellipsis).
// A maxLen of zero or negative is ignored for that field.
func New(fields map[string]int, opts ...Option) (*Truncator, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("truncate: at least one field must be specified")
	}
	t := &Truncator{
		fields:   make(map[string]int, len(fields)),
		ellipsis: defaultEllipsis,
	}
	for _, o := range opts {
		o(t)
	}
	for k, v := range fields {
		if v > 0 {
			t.fields[k] = v
		}
	}
	if len(t.fields) == 0 {
		return nil, fmt.Errorf("truncate: all provided field limits are non-positive")
	}
	return t, nil
}

// Apply truncates configured string fields in the given parsed log entry.
// Non-string fields and absent fields are left untouched.
// The modified map is returned (same map, mutated in place).
func (t *Truncator) Apply(entry map[string]any) map[string]any {
	for field, maxLen := range t.fields {
		v, ok := entry[field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if len(s) > maxLen {
			entry[field] = s[:maxLen] + t.ellipsis
		}
	}
	return entry
}

// ApplyRaw truncates configured string fields in a raw JSON line.
// Returns the modified JSON bytes, or an error if the input is not a JSON object.
func (t *Truncator) ApplyRaw(line []byte) ([]byte, error) {
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return nil, fmt.Errorf("truncate: unmarshal: %w", err)
	}
	t.Apply(entry)
	out, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("truncate: marshal: %w", err)
	}
	return out, nil
}
