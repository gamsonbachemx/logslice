package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/sampler"
)

func buildSamplerLine(i int) string {
	m := map[string]interface{}{
		"ts":  time.Now().UTC().Format(time.RFC3339),
		"msg": fmt.Sprintf("event %d", i),
		"lvl": "info",
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func TestSamplerKeepsSubset(t *testing.T) {
	const total = 200
	var lines []string
	for i := 0; i < total; i++ {
		lines = append(lines, buildSamplerLine(i))
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for _, l := range lines {
			fmt.Fprintln(w, l)
		}
	}))
	defer ts.Close()

	s, err := sampler.New(0.3, 7)
	if err != nil {
		t.Fatalf("sampler.New: %v", err)
	}

	kept := 0
	for _, line := range lines {
		_ = line
		if s.Keep() {
			kept++
		}
	}

	ratio := float64(kept) / float64(total)
	if ratio < 0.15 || ratio > 0.45 {
		t.Errorf("expected ~30%% of %d lines kept, got %d (%.1f%%)", total, kept, ratio*100)
	}
}

func TestSamplerRateOneKeepsAll(t *testing.T) {
	const total = 50
	s, err := sampler.New(1.0, 0)
	if err != nil {
		t.Fatalf("sampler.New: %v", err)
	}

	lines := make([]string, total)
	for i := range lines {
		lines[i] = buildSamplerLine(i)
	}

	kept := 0
	for range lines {
		if s.Keep() {
			kept++
		}
	}
	if kept != total {
		t.Errorf("rate=1.0 should keep all %d lines, kept %d", total, kept)
	}
	_ = strings.Join(lines, "\n")
}
