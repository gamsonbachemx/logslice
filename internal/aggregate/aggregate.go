// Package aggregate provides field-based counting and grouping of log lines.
package aggregate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"
)

// Counter tracks occurrence counts keyed by a field value.
type Counter struct {
	mu    sync.Mutex
	field string
	counts map[string]int
}

// New creates a Counter that groups log lines by the given field name.
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add parses a JSON log line and increments the count for the field value.
// Lines that are not valid JSON or lack the field are counted under "<unknown>".
func (c *Counter) Add(line []byte) {
	var record map[string]interface{}
	key := "<unknown>"
	if err := json.Unmarshal(line, &record); err == nil {
		if v, ok := record[c.field]; ok {
			key = fmt.Sprintf("%v", v)
		}
	}
	c.mu.Lock()
	c.counts[key]++
	c.mu.Unlock()
}

// Result holds a single aggregation bucket.
type Result struct {
	Value string
	Count int
}

// Results returns a sorted slice of Results, descending by count.
func (c *Counter) Results() []Result {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]Result, 0, len(c.counts))
	for k, v := range c.counts {
		out = append(out, Result{Value: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Value < out[j].Value
	})
	return out
}

// Print writes the aggregation table to w.
func (c *Counter) Print(w io.Writer) {
	results := c.Results()
	fmt.Fprintf(w, "%-40s  %s\n", c.field, "count")
	fmt.Fprintf(w, "%-40s  %s\n", "----------------------------------------", "-----")
	for _, r := range results {
		fmt.Fprintf(w, "%-40s  %d\n", r.Value, r.Count)
	}
}
