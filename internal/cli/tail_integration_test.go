package cli_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func TestTailIntegration(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		w.Header().Set("Content-Type", "application/x-ndjson")
		entry := map[string]interface{}{
			"msg":   fmt.Sprintf("log line %d", n),
			"level": "info",
		}
		b, _ := json.Marshal(entry)
		fmt.Fprintln(w, string(b))
	}))
	defer server.Close()

	fetch := func(ctx context.Context) ([]string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var lines []string
		buf := new(strings.Builder)
		b := make([]byte, 512)
		for {
			n, e := resp.Body.Read(b)
			if n > 0 {
				buf.Write(b[:n])
			}
			if e != nil {
				break
			}
		}
		for _, l := range strings.Split(buf.String(), "\n") {
			if strings.TrimSpace(l) != "" {
				lines = append(lines, l)
			}
		}
		return lines, nil
	}

	cfg := tail.Config{PollInterval: 20 * time.Millisecond, MaxRetries: 3}
	tr := tail.New(cfg, fetch)
	out := make(chan string, 20)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = tr.Run(ctx, out)
	close(out)

	var received []string
	for l := range out {
		received = append(received, l)
	}
	if len(received) == 0 {
		t.Error("expected at least one line from tail integration")
	}
}
