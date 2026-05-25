// Package aggregate provides field-based counting and grouping for structured
// JSON log streams.
//
// Use New to create a Counter for a named field, then call Add for each log
// line. Results returns buckets sorted by frequency, and Print renders a
// human-readable summary table.
package aggregate
