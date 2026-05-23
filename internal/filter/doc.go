// Package filter provides structured JSON log parsing and filtering
// capabilities for logslice.
//
// It supports filtering log entries by:
//   - Regular expression patterns matched against the raw JSON line
//   - Time range bounds (start/end) parsed from a configurable timestamp field
//
// Typical usage:
//
//	entry, err := filter.ParseLine(line, "timestamp")
//	if err != nil {
//	    // skip or handle unparseable lines
//	}
//	if filter.Match(entry, opts) {
//	    fmt.Println(entry.Line)
//	}
package filter
