// Package stats tracks streaming and filtering metrics for a logslice session.
package stats

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counter holds atomic counters for a streaming session.
type Counter struct {
	Received  atomic.Int64
	Matched   atomic.Int64
	Skipped   atomic.Int64
	ParseErrs atomic.Int64
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{}
}

// IncReceived increments the total lines received.
func (c *Counter) IncReceived() { c.Received.Add(1) }

// IncMatched increments the matched-lines counter.
func (c *Counter) IncMatched() { c.Matched.Add(1) }

// IncSkipped increments the skipped-lines counter.
func (c *Counter) IncSkipped() { c.Skipped.Add(1) }

// IncParseErr increments the parse-error counter.
func (c *Counter) IncParseErr() { c.ParseErrs.Add(1) }

// Summary returns a human-readable summary string.
func (c *Counter) Summary() string {
	return fmt.Sprintf(
		"received=%d matched=%d skipped=%d parse_errors=%d",
		c.Received.Load(),
		c.Matched.Load(),
		c.Skipped.Load(),
		c.ParseErrs.Load(),
	)
}

// Print writes the summary to w.
func (c *Counter) Print(w io.Writer) {
	fmt.Fprintln(w, c.Summary())
}
