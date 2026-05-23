// Package output handles formatting and writing of filtered log entries.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Format represents the output format for log entries.
type Format string

const (
	// FormatJSON outputs raw JSON lines.
	FormatJSON Format = "json"
	// FormatPretty outputs human-readable formatted lines.
	FormatPretty Format = "pretty"
)

// Writer writes log entries to an output destination.
type Writer struct {
	dest   io.Writer
	format Format
}

// NewWriter creates a new Writer with the given destination and format.
// If dest is nil, os.Stdout is used.
func NewWriter(dest io.Writer, format Format) *Writer {
	if dest == nil {
		dest = os.Stdout
	}
	if format == "" {
		format = FormatJSON
	}
	return &Writer{dest: dest, format: format}
}

// WriteLine writes a single parsed log entry to the destination.
func (w *Writer) WriteLine(entry map[string]interface{}) error {
	switch w.format {
	case FormatPretty:
		return w.writePretty(entry)
	default:
		return w.writeJSON(entry)
	}
}

func (w *Writer) writeJSON(entry map[string]interface{}) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintf(w.dest, "%s\n", data)
	return err
}

func (w *Writer) writePretty(entry map[string]interface{}) error {
	timestamp := ""
	if ts, ok := entry["time"]; ok {
		if tsStr, ok := ts.(string); ok {
			if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
				timestamp = t.Format("2006-01-02 15:04:05")
			}
		}
	}
	level := ""
	if l, ok := entry["level"]; ok {
		level = fmt.Sprintf("[%v]", l)
	}
	message := ""
	if m, ok := entry["message"]; ok {
		message = fmt.Sprintf("%v", m)
	} else if m, ok := entry["msg"]; ok {
		message = fmt.Sprintf("%v", m)
	}
	_, err := fmt.Fprintf(w.dest, "%s %s %s\n", timestamp, level, message)
	return err
}
