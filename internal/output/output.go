// Package output handles writing filtered log lines to an io.Writer.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/highlight"
)

// Format controls how log lines are rendered.
type Format int

const (
	FormatJSON   Format = iota // raw JSON, one line per record
	FormatPretty               // human-readable key=value pairs
)

// Writer writes log lines to an underlying io.Writer.
type Writer struct {
	out         io.Writer
	format      Format
	highlighter *highlight.Highlighter
}

// Options configures a Writer.
type Options struct {
	Format      Format
	Pattern     string // regex to highlight in output (empty = no highlight)
	ColorOutput bool   // enable ANSI colors
}

// NewWriter returns a Writer with the given options.
// If out is nil, os.Stdout is used.
func NewWriter(out io.Writer, opts Options) (*Writer, error) {
	if out == nil {
		out = os.Stdout
	}
	h, err := highlight.New(opts.Pattern, highlight.Cyan, opts.ColorOutput)
	if err != nil {
		return nil, fmt.Errorf("highlight pattern: %w", err)
	}
	return &Writer{out: out, format: opts.Format, highlighter: h}, nil
}

// WriteLine writes a single parsed log record to the underlying writer.
// raw is the original JSON bytes; fields is the decoded map.
func (w *Writer) WriteLine(raw []byte, fields map[string]any) error {
	var line string
	switch w.format {
	case FormatPretty:
		line = w.pretty(fields)
	default:
		line = string(raw)
	}
	line = w.highlighter.Apply(line)
	_, err := fmt.Fprintln(w.out, line)
	return err
}

// pretty renders fields as a space-separated key=value string.
// The "msg" or "message" field is printed first when present.
func (w *Writer) pretty(fields map[string]any) string {
	var sb strings.Builder
	for _, key := range []string{"msg", "message"} {
		if v, ok := fields[key]; ok {
			fmt.Fprintf(&sb, "%s ", v)
			break
		}
	}
	for k, v := range fields {
		if k == "msg" || k == "message" {
			continue
		}
		b, _ := json.Marshal(v)
		fmt.Fprintf(&sb, "%s=%s ", k, string(b))
	}
	return strings.TrimSpace(sb.String())
}
