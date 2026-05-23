package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/output"
)

func TestWriteLineJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)

	entry := map[string]interface{}{
		"time":    "2024-01-15T10:00:00Z",
		"level":   "info",
		"message": "hello world",
	}

	if err := w.WriteLine(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "}") {
		t.Errorf("expected valid JSON line, got: %s", got)
	}
}

func TestWriteLinePretty(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatPretty)

	entry := map[string]interface{}{
		"time":    "2024-01-15T10:00:00Z",
		"level":   "error",
		"message": "something failed",
	}

	if err := w.WriteLine(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[error]") {
		t.Errorf("expected output to contain '[error]', got: %s", got)
	}
	if !strings.Contains(got, "something failed") {
		t.Errorf("expected output to contain message, got: %s", got)
	}
}

func TestNewWriterDefaults(t *testing.T) {
	w := output.NewWriter(nil, "")
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriteLinePrettyMsgField(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatPretty)

	entry := map[string]interface{}{
		"time":  "2024-01-15T10:00:00Z",
		"level": "warn",
		"msg":   "using msg field",
	}

	if err := w.WriteLine(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "using msg field") {
		t.Errorf("expected output to contain msg value, got: %s", got)
	}
}
