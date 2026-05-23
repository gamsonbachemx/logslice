// Package stats provides lightweight atomic counters for tracking
// logslice streaming sessions.
//
// A Counter records how many log lines were received from the remote
// server, how many matched the active filters, how many were skipped,
// and how many could not be parsed as JSON.
//
// Counters are safe for concurrent use via sync/atomic.
package stats
