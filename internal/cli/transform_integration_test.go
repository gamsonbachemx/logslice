package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/logslice/internal/transform"
)

func TestFieldSelectorIntegration(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"started","caller":"main.go:1"}`,
		`{"level":"error","msg":"failed","caller":"main.go:9"}`,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, l := range lines {
			_, _ = w.Write([]byte(l + "\n"))
		}
	}))
	defer ts.Close()

	fs, err := transform.NewFieldSelector("level,msg")
	if err != nil {
		t.Fatalf("selector: %v", err)
	}

	var buf bytes.Buffer
	for _, raw := range lines {
		out, err := fs.ApplyJSON([]byte(raw))
		if err != nil {
			t.Fatalf("ApplyJSON: %v", err)
		}
		buf.Write(out)
		buf.WriteByte('\n')
	}

	result := buf.String()
	if strings.Contains(result, "caller") {
		t.Error("caller field should be absent after projection")
	}

	for _, line := range strings.Split(strings.TrimSpace(result), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if _, ok := m["level"]; !ok {
			t.Error("level field missing")
		}
		if _, ok := m["msg"]; !ok {
			t.Error("msg field missing")
		}
	}
	_ = ts.URL
}
