package config_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
)

func TestValidateValidConfig(t *testing.T) {
	cfg := &config.Config{
		URL:     "https://example.com/logs",
		Timeout: 10 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateMissingURL(t *testing.T) {
	cfg := &config.Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestValidateInvalidURLScheme(t *testing.T) {
	cfg := &config.Config{URL: "ftp://example.com/logs"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for non-http scheme")
	}
}

func TestValidateUntilBeforeSince(t *testing.T) {
	now := time.Now()
	cfg := &config.Config{
		URL:   "https://example.com/logs",
		Since: now,
		Until: now.Add(-time.Hour),
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when until is before since")
	}
}

func TestValidateNegativeTimeout(t *testing.T) {
	cfg := &config.Config{
		URL:     "https://example.com/logs",
		Timeout: -1 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative timeout")
	}
}

func TestValidateNegativeMaxLines(t *testing.T) {
	cfg := &config.Config{
		URL:      "https://example.com/logs",
		MaxLines: -5,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative max-lines")
	}
}

func TestValidateDefaultTimeout(t *testing.T) {
	cfg := &config.Config{URL: "http://example.com/logs"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Timeout != config.DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", config.DefaultTimeout, cfg.Timeout)
	}
}
