// Package levelfilter provides log-level based filtering for structured JSON log lines.
package levelfilter

import (
	"encoding/json"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelFatal,
	"crit":  LevelFatal,
	"critical": LevelFatal,
}

// ParseLevel converts a string to a Level, returning LevelUnknown if unrecognised.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// Filter holds the minimum log level required for a line to pass.
type Filter struct {
	min      Level
	levelKey string
}

// New creates a Filter that passes lines whose level is >= minLevel.
// levelKey is the JSON field name that contains the level string (e.g. "level").
func New(minLevel Level, levelKey string) *Filter {
	if levelKey == "" {
		levelKey = "level"
	}
	return &Filter{min: minLevel, levelKey: levelKey}
}

// Allow returns true when the JSON line's level field meets or exceeds the minimum.
// Lines that cannot be parsed or lack a level field are allowed through.
func (f *Filter) Allow(line []byte) bool {
	var rec map[string]json.RawMessage
	if err := json.Unmarshal(line, &rec); err != nil {
		return true
	}
	raw, ok := rec[f.levelKey]
	if !ok {
		return true
	}
	var val string
	if err := json.Unmarshal(raw, &val); err != nil {
		return true
	}
	l := ParseLevel(val)
	if l == LevelUnknown {
		return true
	}
	return l >= f.min
}
