package transform

import (
	"encoding/json"
	"testing"
)

func TestNewFieldSelectorEmpty(t *testing.T) {
	fs, err := NewFieldSelector("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs.fields) != 0 {
		t.Errorf("expected no fields, got %v", fs.fields)
	}
}

func TestNewFieldSelectorSimple(t *testing.T) {
	fs, err := NewFieldSelector("level,msg,ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs.fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fs.fields))
	}
}

func TestNewFieldSelectorAlias(t *testing.T) {
	fs, err := NewFieldSelector("level:severity,msg:message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.rename["level"] != "severity" {
		t.Errorf("expected alias severity, got %q", fs.rename["level"])
	}
}

func TestNewFieldSelectorInvalidAlias(t *testing.T) {
	_, err := NewFieldSelector("level:")
	if err == nil {
		t.Error("expected error for empty alias, got nil")
	}
}

func TestApplyProjectsFields(t *testing.T) {
	fs, _ := NewFieldSelector("level,msg")
	entry := map[string]any{"level": "info", "msg": "hello", "ts": "2024-01-01"}
	out := fs.Apply(entry)
	if _, ok := out["ts"]; ok {
		t.Error("ts should have been excluded")
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestApplyRenamesFields(t *testing.T) {
	fs, _ := NewFieldSelector("level:severity")
	entry := map[string]any{"level": "warn"}
	out := fs.Apply(entry)
	if _, ok := out["level"]; ok {
		t.Error("original key should be absent after rename")
	}
	if out["severity"] != "warn" {
		t.Errorf("expected severity=warn, got %v", out["severity"])
	}
}

func TestApplyNoFields(t *testing.T) {
	fs, _ := NewFieldSelector("")
	entry := map[string]any{"level": "info", "msg": "hi"}
	out := fs.Apply(entry)
	if len(out) != len(entry) {
		t.Errorf("expected passthrough, got %v", out)
	}
}

func TestApplyJSONRoundtrip(t *testing.T) {
	fs, _ := NewFieldSelector("level,msg")
	raw := []byte(`{"level":"error","msg":"oops","caller":"main.go:10"}`)
	out, err := fs.ApplyJSON(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := result["caller"]; ok {
		t.Error("caller should be excluded")
	}
}

func TestApplyJSONInvalid(t *testing.T) {
	fs, _ := NewFieldSelector("level")
	_, err := fs.ApplyJSON([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
