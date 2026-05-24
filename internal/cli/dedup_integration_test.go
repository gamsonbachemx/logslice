package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/cache"
)

// buildLine returns a minimal JSON log line with the given message.
func buildLine(msg string) string {
	b, _ := json.Marshal(map[string]string{"level": "info", "msg": msg})
	return string(b)
}

func TestDedupCacheFiltersRepeatedLines(t *testing.T) {
	line := buildLine("duplicate event")
	body := strings.Repeat(line+"\n", 5)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	c := cache.New(time.Second, 100)

	seen := 0
	for _, raw := range strings.Split(strings.TrimRight(body, "\n"), "\n") {
		if !c.Seen(raw) {
			seen++
		}
	}

	if seen != 1 {
		t.Fatalf("expected 1 unique line, got %d", seen)
	}
}

func TestDedupCacheAllowsDistinctLines(t *testing.T) {
	c := cache.New(time.Second, 100)

	lines := []string{
		buildLine("event-a"),
		buildLine("event-b"),
		buildLine("event-c"),
	}

	seen := 0
	for _, l := range lines {
		if !c.Seen(l) {
			seen++
		}
	}

	if seen != len(lines) {
		t.Fatalf("expected %d distinct lines, got %d", len(lines), seen)
	}
}

func TestDedupCacheExpiry(t *testing.T) {
	c := cache.New(20*time.Millisecond, 100)

	line := buildLine("expiring event")
	c.Seen(line) // record

	time.Sleep(40 * time.Millisecond)

	if c.Seen(line) {
		t.Fatal("expected line to be unseen after TTL expiry")
	}
}
