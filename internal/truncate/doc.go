// Package truncate provides field-level value truncation for structured JSON
// log entries. It is intended to be used in the logslice pipeline to prevent
// excessively long field values from overwhelming output, while preserving the
// overall structure of the log line.
//
// Usage:
//
//	tr, err := truncate.New(map[string]int{"msg": 200, "error": 500})
//	if err != nil { ... }
//	modified := tr.Apply(entry)
package truncate
