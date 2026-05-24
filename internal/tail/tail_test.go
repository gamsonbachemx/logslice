package tail_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func TestDefaultConfig(t *testing.T) {
	cfg := tail.DefaultConfig()
	if cfg.PollInterval != 2*time.Second {
		t.Errorf("expected 2s poll interval, got %v", cfg.PollInterval)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("expected 5 max retries, got %d", cfg.MaxRetries)
	}
}

func TestRunDeliversLines(t *testing.T) {
	calls := 0
	fetch := func(ctx context.Context) ([]string, error) {
		calls++
		if calls == 1 {
			return []string{"line1", "line2"}, nil
		}
		return nil, nil
	}

	cfg := tail.Config{PollInterval: 10 * time.Millisecond, MaxRetries: 3}
	tr := tail.New(cfg, fetch)
	out := make(chan string, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = tr.Run(ctx, out)
	close(out)

	var got []string
	for l := range out {
		got = append(got, l)
	}
	if len(got) < 2 {
		t.Errorf("expected at least 2 lines, got %d", len(got))
	}
}

func TestRunExceedsMaxRetries(t *testing.T) {
	fetchErr := errors.New("network error")
	fetch := func(ctx context.Context) ([]string, error) {
		return nil, fetchErr
	}

	cfg := tail.Config{PollInterval: 5 * time.Millisecond, MaxRetries: 3}
	tr := tail.New(cfg, fetch)
	out := make(chan string, 10)

	ctx := context.Background()
	err := tr.Run(ctx, out)
	if !errors.Is(err, fetchErr) {
		t.Errorf("expected fetchErr, got %v", err)
	}
}

func TestRunContextCancelled(t *testing.T) {
	fetch := func(ctx context.Context) ([]string, error) {
		return nil, nil
	}

	cfg := tail.Config{PollInterval: 10 * time.Millisecond, MaxRetries: 3}
	tr := tail.New(cfg, fetch)
	out := make(chan string, 10)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := tr.Run(ctx, out)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
