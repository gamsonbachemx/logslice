package filter

import (
	"regexp"
	"testing"
	"time"
)

func TestParseLineValid(t *testing.T) {
	line := `{"level":"info","msg":"started","time":"2024-01-15T10:00:00Z"}`
	entry, err := ParseLine(line, "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected timestamp to be parsed")
	}
	if entry.Raw["level"] != "info" {
		t.Errorf("expected level=info, got %v", entry.Raw["level"])
	}
}

func TestParseLineInvalidJSON(t *testing.T) {
	_, err := ParseLine("not json", "time")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestMatchPattern(t *testing.T) {
	entry := &LogEntry{Line: `{"msg":"connection refused"}`}
	opts := Options{Pattern: regexp.MustCompile(`refused`)}
	if !Match(entry, opts) {
		t.Error("expected match")
	}
	opts.Pattern = regexp.MustCompile(`timeout`)
	if Match(entry, opts) {
		t.Error("expected no match")
	}
}

func TestMatchTimeRange(t *testing.T) {
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entry := &LogEntry{Line: `{}`, Timestamp: ts}

	start := ts.Add(-time.Hour)
	end := ts.Add(time.Hour)
	opts := Options{TimeStart: &start, TimeEnd: &end}
	if !Match(entry, opts) {
		t.Error("expected entry within range to match")
	}

	before := ts.Add(-2 * time.Hour)
	opts = Options{TimeStart: &end, TimeEnd: &before}
	if Match(entry, opts) {
		t.Error("expected entry outside range to not match")
	}
}

func TestMatchNoFilters(t *testing.T) {
	entry := &LogEntry{Line: `{"msg":"hello"}`}
	if !Match(entry, Options{}) {
		t.Error("expected entry with no filters to match")
	}
}
