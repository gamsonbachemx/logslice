package cli

import (
	"strings"
	"testing"
	"time"
)

func TestParseFlagsAllOptions(t *testing.T) {
	cfg, err := parseFlags([]string{
		"--url", "http://logs.example.com",
		"--pattern", "panic",
		"--since", "2024-01-01T00:00:00Z",
		"--until", "2024-12-31T23:59:59Z",
		"--pretty",
		"--timeout", "10s",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "http://logs.example.com" {
		t.Errorf("unexpected URL: %s", cfg.URL)
	}
	if cfg.Pattern != "panic" {
		t.Errorf("unexpected pattern: %s", cfg.Pattern)
	}
	if cfg.Since != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected since: %s", cfg.Since)
	}
	if cfg.Until != "2024-12-31T23:59:59Z" {
		t.Errorf("unexpected until: %s", cfg.Until)
	}
	if !cfg.Pretty {
		t.Error("expected pretty=true")
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.Timeout)
	}
}

func TestParseFlagsEmptyURL(t *testing.T) {
	_, err := parseFlags([]string{"--url", ""})
	if err == nil || !strings.Contains(err.Error(), "--url is required") {
		t.Fatalf("expected url-required error, got %v", err)
	}
}

func TestParseFlagsNoArgs(t *testing.T) {
	_, err := parseFlags([]string{})
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestParseFlagsPrettyDefault(t *testing.T) {
	cfg, err := parseFlags([]string{"--url", "http://example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pretty {
		t.Error("expected pretty to default to false")
	}
}
