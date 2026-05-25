// Package multiurl provides fan-out streaming across multiple log source URLs.
//
// Fan accepts a slice of URLs and a Streamer function, launching one goroutine
// per URL and merging all resulting lines into a single output channel. Each
// line is tagged with its source URL so consumers can distinguish origins.
//
// Example:
//
//	lineCh, errCh := multiurl.Fan(ctx, urls, func(ctx context.Context, url string, out chan<- multiurl.Line) error {
//		return fetch.Stream(ctx, url, func(line string) { out <- multiurl.Line{Source: url, Data: line} })
//	})
package multiurl
