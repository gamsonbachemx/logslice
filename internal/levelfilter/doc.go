// Package levelfilter implements minimum log-level filtering for JSON log lines.
//
// It parses a configurable level field (default: "level") from each JSON record
// and compares it against a minimum severity threshold. Lines at or above the
// threshold are allowed; lines below are dropped. Lines that cannot be parsed
// or that lack a level field are passed through unchanged so that non-standard
// records are never silently discarded.
//
// Supported level strings (case-insensitive):
//
//	debug < info < warn/warning < error/err < fatal/crit/critical
package levelfilter
