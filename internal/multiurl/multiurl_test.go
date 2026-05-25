package multiurl_test

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/multiurl"
)

func TestFanMergesAllSources(t *testing.T) {
	urls := []string{"http://host1", "http://host2", "http://host3"}

	streamer := func(ctx context.Context, url string, out chan<- multiurl.Line) error {
		for i := 0; i < 3; i++ {
			out <- multiurl.Line{Source: url, Data: fmt.Sprintf("%s-line%d", url, i)}
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lineCh, errCh := multiurl.Fan(ctx, urls, streamer)

	var mu sync.Mutex
	var received []string
	for l := range lineCh {
		mu.Lock()
		received = append(received, l.Data)
		mu.Unlock()
	}
	for err := range errCh {
		t.Errorf("unexpected error: %v", err)
	}

	if len(received) != 9 {
		t.Fatalf("expected 9 lines, got %d", len(received))
	}
}

func TestFanPropagatesErrors(t *testing.T) {
	urls := []string{"http://bad"}

	streamer := func(ctx context.Context, url string, out chan<- multiurl.Line) error {
		return fmt.Errorf("connection refused")
	}

	ctx := context.Background()
	_, errCh := multiurl.Fan(ctx, urls, streamer)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestFanSourceTagging(t *testing.T) {
	urls := []string{"http://alpha", "http://beta"}

	streamer := func(ctx context.Context, url string, out chan<- multiurl.Line) error {
		out <- multiurl.Line{Source: url, Data: "msg"}
		return nil
	}

	ctx := context.Background()
	lineCh, _ := multiurl.Fan(ctx, urls, streamer)

	sources := map[string]bool{}
	for l := range lineCh {
		sources[l.Source] = true
	}

	keys := []string{}
	for k := range sources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(keys) != 2 || keys[0] != "http://alpha" || keys[1] != "http://beta" {
		t.Errorf("unexpected sources: %v", keys)
	}
}

func TestFanEmptyURLs(t *testing.T) {
	streamer := func(ctx context.Context, url string, out chan<- multiurl.Line) error {
		return nil
	}

	ctx := context.Background()
	lineCh, errCh := multiurl.Fan(ctx, []string{}, streamer)

	for range lineCh {
		t.Error("expected no lines")
	}
	for range errCh {
		t.Error("expected no errors")
	}
}
