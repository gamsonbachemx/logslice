// Package cli parses command-line flags and wires together fetch, filter,
// and output to stream structured JSON logs from a remote server.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/fetch"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
)

// Config holds all runtime options resolved from CLI flags.
type Config struct {
	URL       string
	Pattern   string
	Since     string
	Until     string
	Pretty    bool
	Timeout   time.Duration
	Out       io.Writer
}

// Run parses args, validates them, and executes the main streaming pipeline.
func Run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}
	return stream(cfg)
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	url := fs.String("url", "", "URL of the remote log endpoint (required)")
	pattern := fs.String("pattern", "", "Regex pattern to filter log messages")
	since := fs.String("since", "", "Include logs after this time (RFC3339)")
	until := fs.String("until", "", "Include logs before this time (RFC3339)")
	pretty := fs.Bool("pretty", false, "Pretty-print log output")
	timeout := fs.Duration("timeout", 30*time.Second, "HTTP request timeout")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *url == "" {
		return nil, fmt.Errorf("--url is required")
	}

	return &Config{
		URL:     *url,
		Pattern: *pattern,
		Since:   *since,
		Until:   *until,
		Pretty:  *pretty,
		Timeout: *timeout,
		Out:     os.Stdout,
	}, nil
}

func stream(cfg *Config) error {
	client := fetch.NewClient(cfg.Timeout)

	lines, err := client.Stream(cfg.URL)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	w, err := output.NewWriter(cfg.Out, cfg.Pretty)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}

	for line := range lines {
		entry, parseErr := filter.ParseLine(line)
		if parseErr != nil {
			continue
		}
		if !filter.Match(entry, cfg.Pattern, cfg.Since, cfg.Until) {
			continue
		}
		if writeErr := w.WriteLine(entry); writeErr != nil {
			return fmt.Errorf("write: %w", writeErr)
		}
	}

	return nil
}
