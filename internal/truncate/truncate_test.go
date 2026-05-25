package truncate_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logslice/internal/truncate"
)

func TestNewNoFields(t *testing.T) {
	_, err := truncate.New(nil)
	if err == nil {
		t.Fatal("expected error for nil fields")
	}
}

func TestNewAllNonPositive(t *testing.T) {
	_, err := truncate.New(map[string]int{"msg": 0, "level": -1})
	if err == nil {
		t.Fatal("expected error when all limits are non-positive")
	}
}

func TestApplyTruncatesLongField(t *testing.T) {
	tr, err := truncate.New(map[string]int{"msg": 10})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	entry := map[string]any{"msg": "hello world, this is long", "level": "info"}
	out := tr.Apply(entry)
	got, _ := out["msg"].(string)
	if len(got) != 10+len("...") {
		t.Errorf("expected length %d, got %d (%q)", 10+3, len(got), got)
	}
	if got != "hello worl..." {
		t.Errorf("unexpected value: %q", got)
	}
}

func TestApplySkipsShortField(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"msg": 50})
	entry := map[string]any{"msg": "short"}
	out := tr.Apply(entry)
	if out["msg"] != "short" {
		t.Errorf("expected unchanged value, got %v", out["msg"])
	}
}

func TestApplySkipsNonStringField(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"code": 3})
	entry := map[string]any{"code": 404}
	out := tr.Apply(entry)
	if out["code"] != 404 {
		t.Errorf("non-string field should be untouched, got %v", out["code"])
	}
}

func TestApplySkipsMissingField(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"msg": 5})
	entry := map[string]any{"level": "warn"}
	out := tr.Apply(entry)
	if _, ok := out["msg"]; ok {
		t.Error("missing field should not be added")
	}
}

func TestApplyRawValid(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"msg": 5})
	line := []byte(`{"msg":"hello world","level":"info"}`)
	out, err := tr.ApplyRaw(line)
	if err != nil {
		t.Fatalf("ApplyRaw error: %v", err)
	}
	var entry map[string]any
	if err := json.Unmarshal(out, &entry); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if entry["msg"] != "hello..." {
		t.Errorf("unexpected msg: %v", entry["msg"])
	}
}

func TestApplyRawInvalidJSON(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"msg": 5})
	_, err := tr.ApplyRaw([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWithEllipsisOption(t *testing.T) {
	tr, _ := truncate.New(map[string]int{"msg": 4}, truncate.WithEllipsis("~"))
	entry := map[string]any{"msg": "hello world"}
	out := tr.Apply(entry)
	if out["msg"] != "hell~" {
		t.Errorf("expected 'hell~', got %v", out["msg"])
	}
}
