// Package multiurl provides fan-out streaming across multiple log source URLs.
package multiurl

import (
	"context"
	"sync"
)

// Line represents a single log line along with the source URL it came from.
type Line struct {
	Source string
	Data   string
}

// Streamer is a function that streams lines from a single URL into a channel.
type Streamer func(ctx context.Context, url string, out chan<- Line) error

// Fan streams from all urls concurrently, merging results into a single channel.
// The returned channel is closed once all sources finish or the context is cancelled.
// Errors from individual sources are sent to the errCh channel (non-blocking).
func Fan(ctx context.Context, urls []string, streamer Streamer) (<-chan Line, <-chan error) {
	out := make(chan Line, len(urls)*32)
	errCh := make(chan error, len(urls))

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := streamer(ctx, url, out); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errCh)
	}()

	return out, errCh
}
