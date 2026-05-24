package tail

import (
	"context"
	"time"
)

// Config holds configuration for tail mode.
type Config struct {
	// PollInterval is how often to re-fetch new lines when tailing.
	PollInterval time.Duration
	// MaxRetries is the number of consecutive fetch errors before giving up.
	MaxRetries int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PollInterval: 2 * time.Second,
		MaxRetries:   5,
	}
}

// Tailer repeatedly calls the provided fetch function at the configured
// interval, forwarding new lines to the out channel until ctx is cancelled.
type Tailer struct {
	cfg   Config
	fetch func(ctx context.Context) ([]string, error)
}

// New creates a new Tailer.
func New(cfg Config, fetch func(ctx context.Context) ([]string, error)) *Tailer {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = DefaultConfig().PollInterval
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = DefaultConfig().MaxRetries
	}
	return &Tailer{cfg: cfg, fetch: fetch}
}

// Run starts tailing and sends lines to out. It returns when ctx is done
// or consecutive errors exceed MaxRetries.
func (t *Tailer) Run(ctx context.Context, out chan<- string) error {
	consecErrors := 0
	ticker := time.NewTicker(t.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			lines, err := t.fetch(ctx)
			if err != nil {
				consecErrors++
				if consecErrors >= t.cfg.MaxRetries {
					return err
				}
				continue
			}
			consecErrors = 0
			for _, line := range lines {
				select {
				case out <- line:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}
