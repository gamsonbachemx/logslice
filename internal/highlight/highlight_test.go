package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestNewNoPattern(t *testing.T) {
	h, err := highlight.New("", highlight.Cyan, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "hello world"
	if got := h.Apply(input); got != input {
		t.Errorf("expected unchanged input, got %q", got)
	}
}

func TestNewDisabled(t *testing.T) {
	h, err := highlight.New("error", highlight.Red, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "an error occurred"
	if got := h.Apply(input); got != input {
		t.Errorf("expected unchanged input when disabled, got %q", got)
	}
}

func TestNewInvalidPattern(t *testing.T) {
	_, err := highlight.New("[invalid", highlight.Red, true)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestApplyHighlightsMatch(t *testing.T) {
	h, err := highlight.New("error", highlight.Red, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "an error occurred"
	got := h.Apply(input)
	if !strings.Contains(got, highlight.Red) {
		t.Errorf("expected ANSI red code in output, got %q", got)
	}
	if !strings.Contains(got, highlight.Reset) {
		t.Errorf("expected ANSI reset code in output, got %q", got)
	}
	stripped := highlight.StripANSI(got)
	if stripped != input {
		t.Errorf("stripped output mismatch: want %q, got %q", input, stripped)
	}
}

func TestApplyMultipleMatches(t *testing.T) {
	h, err := highlight.New(`\d+`, highlight.Yellow, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "line 42 and item 7"
	got := h.Apply(input)
	count := strings.Count(got, highlight.Yellow)
	if count != 2 {
		t.Errorf("expected 2 highlights, got %d in %q", count, got)
	}
}

func TestStripANSI(t *testing.T) {
	colored := highlight.Red + "hello" + highlight.Reset
	if got := highlight.StripANSI(colored); got != "hello" {
		t.Errorf("StripANSI: want %q, got %q", "hello", got)
	}
}
