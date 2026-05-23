package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

func TestNewCounter(t *testing.T) {
	c := stats.New()
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestIncReceivedAndMatched(t *testing.T) {
	c := stats.New()
	c.IncReceived()
	c.IncReceived()
	c.IncMatched()

	if got := c.Received.Load(); got != 2 {
		t.Errorf("Received: want 2, got %d", got)
	}
	if got := c.Matched.Load(); got != 1 {
		t.Errorf("Matched: want 1, got %d", got)
	}
}

func TestIncSkippedAndParseErr(t *testing.T) {
	c := stats.New()
	c.IncSkipped()
	c.IncParseErr()
	c.IncParseErr()

	if got := c.Skipped.Load(); got != 1 {
		t.Errorf("Skipped: want 1, got %d", got)
	}
	if got := c.ParseErrs.Load(); got != 2 {
		t.Errorf("ParseErrs: want 2, got %d", got)
	}
}

func TestSummaryFormat(t *testing.T) {
	c := stats.New()
	c.IncReceived()
	c.IncMatched()
	c.IncSkipped()
	c.IncParseErr()

	s := c.Summary()
	for _, want := range []string{"received=1", "matched=1", "skipped=1", "parse_errors=1"} {
		if !strings.Contains(s, want) {
			t.Errorf("Summary() missing %q in %q", want, s)
		}
	}
}

func TestPrint(t *testing.T) {
	c := stats.New()
	c.IncReceived()
	var buf bytes.Buffer
	c.Print(&buf)
	if !strings.Contains(buf.String(), "received=1") {
		t.Errorf("Print() output missing 'received=1': %q", buf.String())
	}
}
