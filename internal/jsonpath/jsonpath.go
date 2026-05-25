// Package jsonpath provides a simple dot-notation extractor for nested
// JSON log fields (e.g. "metadata.request.method").
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extractor resolves one or more dot-notation paths from a raw JSON line.
type Extractor struct {
	paths []string
}

// New returns an Extractor for the given dot-notation paths.
// An error is returned if no paths are provided.
func New(paths []string) (*Extractor, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("jsonpath: at least one path is required")
	}
	return &Extractor{paths: paths}, nil
}

// Extract returns a map of path → value for each configured path.
// Paths that cannot be resolved are omitted from the result.
func (e *Extractor) Extract(line []byte) (map[string]any, error) {
	var root map[string]any
	if err := json.Unmarshal(line, &root); err != nil {
		return nil, fmt.Errorf("jsonpath: unmarshal: %w", err)
	}

	result := make(map[string]any, len(e.paths))
	for _, p := range e.paths {
		if v, ok := resolve(root, p); ok {
			result[p] = v
		}
	}
	return result, nil
}

// Paths returns the configured dot-notation paths.
func (e *Extractor) Paths() []string {
	out := make([]string, len(e.paths))
	copy(out, e.paths)
	return out
}

// resolve walks a nested map following the dot-separated segments of path.
func resolve(node map[string]any, path string) (any, bool) {
	segments := strings.SplitN(path, ".", 2)
	v, ok := node[segments[0]]
	if !ok {
		return nil, false
	}
	if len(segments) == 1 {
		return v, true
	}
	child, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	return resolve(child, segments[1])
}
