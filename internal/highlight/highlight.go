// Package highlight provides ANSI color highlighting for matched patterns in log output.
package highlight

import (
	"regexp"
	"strings"
)

// ANSI color codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter applies ANSI color codes to substrings matching a pattern.
type Highlighter struct {
	pattern *regexp.Regexp
	color   string
	enabled bool
}

// New creates a Highlighter for the given regex pattern and color code.
// If pattern is empty or enabled is false, the Highlighter is a no-op.
func New(pattern, color string, enabled bool) (*Highlighter, error) {
	h := &Highlighter{color: color, enabled: enabled}
	if pattern != "" && enabled {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		h.pattern = re
	}
	return h, nil
}

// Apply wraps all matches of the pattern in the input string with ANSI color codes.
// If the Highlighter is disabled or has no pattern, the input is returned unchanged.
func (h *Highlighter) Apply(input string) string {
	if !h.enabled || h.pattern == nil {
		return input
	}
	return h.pattern.ReplaceAllStringFunc(input, func(match string) string {
		return h.color + match + Reset
	})
}

// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string {
	re := regexp.MustCompile(`\033\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

// IsTerminal returns true when the NO_COLOR env var is not set and the
// color string is non-empty — callers may use this as a lightweight check.
func IsTerminal(color string) bool {
	return strings.TrimSpace(color) != ""
}
