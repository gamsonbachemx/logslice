// Package redact provides field-level redaction for structured log lines.
// It replaces the values of specified JSON fields with a configurable mask string.
package redact

import (
	"encoding/json"
	"fmt"
)

const defaultMask = "[REDACTED]"

// Redactor replaces sensitive field values in parsed log lines.
type Redactor struct {
	fields map[string]struct{}
	mask   string
}

// New creates a Redactor that masks the given field names.
// If mask is empty, the default mask "[REDACTED]" is used.
// Returns an error if no fields are provided.
func New(fields []string, mask string) (*Redactor, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("redact: at least one field name is required")
	}
	if mask == "" {
		mask = defaultMask
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return &Redactor{fields: set, mask: mask}, nil
}

// Apply redacts sensitive fields in the given parsed log map.
// It returns a new map with the specified fields replaced by the mask.
// Fields not present in the map are silently ignored.
func (r *Redactor) Apply(record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(record))
	for k, v := range record {
		if _, ok := r.fields[k]; ok {
			out[k] = r.mask
		} else {
			out[k] = v
		}
	}
	return out
}

// ApplyBytes parses a JSON line, redacts fields, and re-serialises it.
// Returns the redacted JSON bytes or an error if parsing fails.
func (r *Redactor) ApplyBytes(line []byte) ([]byte, error) {
	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("redact: failed to parse JSON: %w", err)
	}
	redacted := r.Apply(record)
	out, err := json.Marshal(redacted)
	if err != nil {
		return nil, fmt.Errorf("redact: failed to serialise redacted record: %w", err)
	}
	return out, nil
}
