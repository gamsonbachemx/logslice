// Package cli implements the command-line interface for logslice.
//
// It exposes a single Run function that accepts os.Args[1:] and orchestrates
// the full pipeline:
//
//  1. Parse flags (--url, --pattern, --since, --until, --pretty, --timeout).
//  2. Open an HTTP stream via the fetch package.
//  3. Parse each newline-delimited JSON log entry via the filter package.
//  4. Apply regex and time-range filters.
//  5. Write matching entries to stdout via the output package.
//
// Example usage:
//
//	logslice --url https://logs.example.com/stream \
//	         --pattern "error" \
//	         --since 2024-01-01T00:00:00Z \
//	         --pretty
package cli
