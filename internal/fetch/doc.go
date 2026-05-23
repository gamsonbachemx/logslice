// Package fetch provides an HTTP client for streaming newline-delimited
// JSON log data from remote servers. It supports configurable timeouts
// and exposes a line-by-line streaming interface suitable for large or
// continuously written log endpoints.
//
// Usage:
//
//	client := fetch.NewClient(30 * time.Second)
//	err := client.Stream(ctx, "https://logs.example.com/stream", func(line []byte) error {
//		// process each line
//		return nil
//	})
package fetch
