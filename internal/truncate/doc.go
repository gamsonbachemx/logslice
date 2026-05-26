// Package truncate provides field-level value truncation for structured JSON
// log entries. It is intended to be used in the logslice pipeline to prevent
// excessively long field values from overwhelming output, while preserving the
// overall structure of the log line.
//
// A Truncator is configured with a map of field names to maximum byte lengths.
// When applied to a log entry, any string field whose value exceeds the
// configured limit is truncated and a suffix (e.g. "...") is appended to
// indicate that the value was shortened.
//
// Fields not present in the configuration are left unchanged. Non-string
// field values are also left unchanged regardless of their size.
//
// Usage:
//
//	tr, err := truncate.New(map[string]int{"msg": 200, "error": 500})
//	if err != nil { ... }
//	modified := tr.Apply(entry)
package truncate
