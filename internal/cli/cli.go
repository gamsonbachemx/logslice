// Package cli wires together flags, fetching, filtering, output, and stats.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/fetch"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/stats"
)

// config holds parsed CLI options.
type config struct {
	URL     string
	Pattern string
	From    string
	To      string
	Pretty  bool
	Timeout time.Duration
}

// Run is the entry point called from main.
func Run(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		return 2
	}
	if cfg.URL == "" {
		fmt.Fprintln(stderr, "error: --url is required")
		return 2
	}

	client := fetch.NewClient(cfg.Timeout)
	writer := output.NewWriter(stdout, cfg.Pretty)
	counters := stats.New()

	ctx := context.Background()
	if err := stream(ctx, client, writer, counters, cfg); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		counters.Print(stderr)
		return 1
	}
	counters.Print(stderr)
	return 0
}

func parseFlags(args []string, stderr io.Writer) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)

	cfg := &config{}
	fs.StringVar(&cfg.URL, "url", "", "HTTP(S) URL of the log stream (required)")
	fs.StringVar(&cfg.Pattern, "pattern", "", "Regex pattern to match against log lines")
	fs.StringVar(&cfg.From, "from", "", "Include lines at or after this RFC3339 timestamp")
	fs.StringVar(&cfg.To, "to", "", "Include lines at or before this RFC3339 timestamp")
	fs.BoolVar(&cfg.Pretty, "pretty", false, "Pretty-print JSON output")
	fs.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "HTTP client timeout")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return cfg, nil
}

func stream(
	ctx context.Context,
	client *fetch.Client,
	writer *output.Writer,
	counters *stats.Counter,
	cfg *config,
) error {
	lines, errCh := client.Stream(ctx, cfg.URL)
	for line := range lines {
		counters.IncReceived()
		entry, err := filter.ParseLine(line)
		if err != nil {
			counters.IncParseErr()
			continue
		}
		if !filter.Match(entry, cfg.Pattern, cfg.From, cfg.To) {
			counters.IncSkipped()
			continue
		}
		counters.IncMatched()
		if err := writer.WriteLine(entry); err != nil {
			return err
		}
	}
	return <-errCh
}

// defaultStdout/Stderr used by main.
var defaultStdout io.Writer = os.Stdout
var defaultStderr io.Writer = os.Stderr
