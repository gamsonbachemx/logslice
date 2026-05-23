// Package fetch provides functionality for retrieving log streams
// from remote servers over SSH or HTTP.
package fetch

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds the configuration for a remote log fetch operation.
type Config struct {
	URL     string
	Timeout time.Duration
	Headers map[string]string
}

// DefaultTimeout is the default HTTP request timeout.
const DefaultTimeout = 30 * time.Second

// Client wraps an HTTP client for fetching remote log streams.
type Client struct {
	http *http.Client
}

// NewClient creates a new fetch Client with the given timeout.
// If timeout is zero, DefaultTimeout is used.
func NewClient(timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &Client{
		http: &http.Client{Timeout: timeout},
	}
}

// Stream fetches the remote resource at cfg.URL and streams each line
// to the returned channel. The channel is closed when the response body
// is fully consumed or an error occurs. Any error is sent on errCh.
func (c *Client) Stream(cfg Config) (<-chan string, <-chan error) {
	lines := make(chan string)
	errCh := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errCh)

		req, err := http.NewRequest(http.MethodGet, cfg.URL, nil)
		if err != nil {
			errCh <- fmt.Errorf("fetch: building request: %w", err)
			return
		}
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("fetch: executing request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("fetch: unexpected status %d for %s", resp.StatusCode, cfg.URL)
			return
		}

		if err := readLines(resp.Body, lines); err != nil {
			errCh <- fmt.Errorf("fetch: reading body: %w", err)
		}
	}()

	return lines, errCh
}

// readLines scans r line-by-line, sending each non-empty line to ch.
func readLines(r io.Reader, ch chan<- string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			ch <- line
		}
	}
	return scanner.Err()
}
