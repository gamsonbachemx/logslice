package cli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/cli"
)

func newTestServer(lines []map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, l := range lines {
			b, _ := json.Marshal(l)
			fmt.Fprintf(w, "%s\n", b)
		}
	}))
}

func TestRunMissingURL(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := cli.Run([]string{}, &out, &errBuf)
	if code != 2 {
		t.Errorf("want exit 2, got %d", code)
	}
	if !strings.Contains(errBuf.String(), "--url") {
		t.Errorf("expected --url mention in stderr: %q", errBuf.String())
	}
}

func TestRunInvalidFlag(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := cli.Run([]string{"--notaflag"}, &out, &errBuf)
	if code != 2 {
		t.Errorf("want exit 2, got %d", code)
	}
}

func TestParseFlagsDefaults(t *testing.T) {
	var out, errBuf bytes.Buffer
	server := newTestServer([]map[string]any{
		{"msg": "hello", "time": "2024-01-01T00:00:00Z"},
	})
	defer server.Close()

	code := cli.Run([]string{"--url", server.URL}, &out, &errBuf)
	if code != 0 {
		t.Errorf("want exit 0, got %d; stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "hello") {
		t.Errorf("expected 'hello' in output: %q", out.String())
	}
}

func TestStreamFiltersAndWrites(t *testing.T) {
	var out, errBuf bytes.Buffer
	server := newTestServer([]map[string]any{
		{"msg": "match this", "time": "2024-06-01T10:00:00Z"},
		{"msg": "ignore me", "time": "2024-06-01T11:00:00Z"},
	})
	defer server.Close()

	code := cli.Run([]string{
		"--url", server.URL,
		"--pattern", "match",
	}, &out, &errBuf)

	if code != 0 {
		t.Errorf("want exit 0, got %d; stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "match this") {
		t.Errorf("expected 'match this' in output: %q", out.String())
	}
	if strings.Contains(out.String(), "ignore me") {
		t.Errorf("did not expect 'ignore me' in output: %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "matched=1") {
		t.Errorf("expected stats in stderr: %q", errBuf.String())
	}
}
