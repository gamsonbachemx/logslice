package jsonpath_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/jsonpath"
)

func TestNewNoPaths(t *testing.T) {
	_, err := jsonpath.New(nil)
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestExtractTopLevel(t *testing.T) {
	ex, err := jsonpath.New([]string{"level", "msg"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	line := []byte(`{"level":"info","msg":"hello","ts":1}`)
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if res["level"] != "info" {
		t.Errorf("level: got %v, want info", res["level"])
	}
	if res["msg"] != "hello" {
		t.Errorf("msg: got %v, want hello", res["msg"])
	}
}

func TestExtractNested(t *testing.T) {
	ex, _ := jsonpath.New([]string{"metadata.request.method"})
	line := []byte(`{"metadata":{"request":{"method":"GET"}}}`)
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if res["metadata.request.method"] != "GET" {
		t.Errorf("nested: got %v, want GET", res["metadata.request.method"])
	}
}

func TestExtractMissingPathOmitted(t *testing.T) {
	ex, _ := jsonpath.New([]string{"does.not.exist", "level"})
	line := []byte(`{"level":"warn"}`)
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	if _, ok := res["does.not.exist"]; ok {
		t.Error("expected missing path to be omitted")
	}
	if res["level"] != "warn" {
		t.Errorf("level: got %v, want warn", res["level"])
	}
}

func TestExtractInvalidJSON(t *testing.T) {
	ex, _ := jsonpath.New([]string{"level"})
	_, err := ex.Extract([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestPathsReturnsConfigured(t *testing.T) {
	paths := []string{"a", "b.c"}
	ex, _ := jsonpath.New(paths)
	got := ex.Paths()
	if len(got) != len(paths) {
		t.Fatalf("Paths len: got %d, want %d", len(got), len(paths))
	}
	for i, p := range paths {
		if got[i] != p {
			t.Errorf("Paths[%d]: got %q, want %q", i, got[i], p)
		}
	}
}
