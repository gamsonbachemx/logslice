package aggregate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/aggregate"
)

func TestNewCounter(t *testing.T) {
	c := aggregate.New("level")
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}

func TestAddValidLines(t *testing.T) {
	c := aggregate.New("level")
	c.Add([]byte(`{"level":"info","msg":"a"}`))
	c.Add([]byte(`{"level":"info","msg":"b"}`))
	c.Add([]byte(`{"level":"error","msg":"c"}`))

	results := c.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(results))
	}
	if results[0].Value != "info" || results[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "error" || results[1].Count != 1 {
		t.Errorf("expected error=1, got %s=%d", results[1].Value, results[1].Count)
	}
}

func TestAddInvalidJSON(t *testing.T) {
	c := aggregate.New("level")
	c.Add([]byte(`not json`))

	results := c.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(results))
	}
	if results[0].Value != "<unknown>" {
		t.Errorf("expected <unknown>, got %s", results[0].Value)
	}
}

func TestAddMissingField(t *testing.T) {
	c := aggregate.New("level")
	c.Add([]byte(`{"msg":"no level here"}`))

	results := c.Results()
	if results[0].Value != "<unknown>" {
		t.Errorf("expected <unknown>, got %s", results[0].Value)
	}
}

func TestResultsSortedByCountDesc(t *testing.T) {
	c := aggregate.New("service")
	for i := 0; i < 5; i++ {
		c.Add([]byte(`{"service":"alpha"}`))
	}
	for i := 0; i < 3; i++ {
		c.Add([]byte(`{"service":"beta"}`))
	}
	c.Add([]byte(`{"service":"gamma"}`))

	results := c.Results()
	if results[0].Value != "alpha" {
		t.Errorf("expected alpha first, got %s", results[0].Value)
	}
	if results[2].Value != "gamma" {
		t.Errorf("expected gamma last, got %s", results[2].Value)
	}
}

func TestPrintOutput(t *testing.T) {
	c := aggregate.New("level")
	c.Add([]byte(`{"level":"warn"}`))
	c.Add([]byte(`{"level":"warn"}`))

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "level") {
		t.Error("expected header to contain field name")
	}
	if !strings.Contains(out, "warn") {
		t.Error("expected output to contain 'warn'")
	}
	if !strings.Contains(out, "2") {
		t.Error("expected count 2 in output")
	}
}
