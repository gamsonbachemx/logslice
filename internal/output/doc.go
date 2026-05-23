// Package output provides formatting and writing utilities for structured
// log entries produced by logslice.
//
// It supports two output formats:
//
//   - FormatJSON: emits each log entry as a compact JSON line, preserving
//     all original fields.
//
//   - FormatPretty: emits a human-readable line containing the timestamp,
//     log level, and message extracted from well-known fields ("time",
//     "level", "message" or "msg").
//
// Example usage:
//
//	w := output.NewWriter(os.Stdout, output.FormatPretty)
//	w.WriteLine(entry)
package output
