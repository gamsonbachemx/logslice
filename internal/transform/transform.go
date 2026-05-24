package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldSelector extracts or renames fields from a parsed log entry.
type FieldSelector struct {
	fields []string
	rename map[string]string
}

// NewFieldSelector creates a FieldSelector from a comma-separated field list.
// Each entry may be "field" or "field:alias".
func NewFieldSelector(spec string) (*FieldSelector, error) {
	if spec == "" {
		return &FieldSelector{}, nil
	}
	fs := &FieldSelector{
		rename: make(map[string]string),
	}
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if idx := strings.Index(part, ":"); idx >= 0 {
			src := strings.TrimSpace(part[:idx])
			dst := strings.TrimSpace(part[idx+1:])
			if src == "" || dst == "" {
				return nil, fmt.Errorf("invalid field alias %q", part)
			}
			fs.fields = append(fs.fields, src)
			fs.rename[src] = dst
		} else {
			fs.fields = append(fs.fields, part)
		}
	}
	return fs, nil
}

// Apply projects the given map to only the selected fields, applying aliases.
// If no fields are configured, the original map is returned unchanged.
func (fs *FieldSelector) Apply(entry map[string]any) map[string]any {
	if len(fs.fields) == 0 {
		return entry
	}
	out := make(map[string]any, len(fs.fields))
	for _, f := range fs.fields {
		val, ok := entry[f]
		if !ok {
			continue
		}
		key := f
		if alias, ok := fs.rename[f]; ok {
			key = alias
		}
		out[key] = val
	}
	return out
}

// ApplyJSON applies the field selector to a raw JSON line and returns the
// re-serialised JSON bytes.
func (fs *FieldSelector) ApplyJSON(line []byte) ([]byte, error) {
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return nil, fmt.Errorf("transform: unmarshal: %w", err)
	}
	projected := fs.Apply(entry)
	out, err := json.Marshal(projected)
	if err != nil {
		return nil, fmt.Errorf("transform: marshal: %w", err)
	}
	return out, nil
}
