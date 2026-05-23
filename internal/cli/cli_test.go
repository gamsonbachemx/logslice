package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer(lines []string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, l := range lines {
			fmt.Fprintln(w, l)
		}
	}))
}

func TestRunMissingURL(t *testing.T) {
	err := Run([]string{})
	if err == nil || !strings.Contains(err.Error(), "--url is required") {
		t.Fatalf("expected url-required error, got %v", err)
	}
}

func TestRunInvalidFlag(t *testing.T) {
	err := Run([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseFlagsDefaults(t *testing.T) {
	cfg, err := parseFlags([]string{"--url", "http://example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pretty {
		t.Error("expected pretty=false by default")
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", cfg.Timeout)
	}
	if cfg.Pattern != "" {
		t.Error("expected empty pattern by default")
	}
}

func TestStreamFiltersAndWrites(t *testing.T) {
	entry := map[string]interface{}{
		"time":    time.Now().UTC().Format(time.RFC3339),
		"level":   "info",
		"message": "hello world",
	}
	raw, _ := json.Marshal(entry)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(raw)
		w.Write([]byte("\n"))
		w.Write([]byte("not json\n"))
	}))
	defer srv.Close()

	var buf bytes.Buffer
	cfg := &Config{
		URL:     srv.URL,
		Pattern: "hello",
		Timeout: 5 * time.Second,
		Out:     &buf,
	}

	if err := stream(cfg); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", buf.String())
	}
}
