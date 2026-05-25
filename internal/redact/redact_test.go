package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logslice/internal/redact"
)

func TestNewNoFields(t *testing.T) {
	_, err := redact.New(nil, "")
	if err == nil {
		t.Fatal("expected error for empty fields, got nil")
	}
}

func TestNewDefaultMask(t *testing.T) {
	r, err := redact.New([]string{"password"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{"password": "secret", "user": "alice"}
	out := r.Apply(record)
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("expected alice, got %v", out["user"])
	}
}

func TestApplyCustomMask(t *testing.T) {
	r, err := redact.New([]string{"token"}, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{"token": "abc123", "level": "info"}
	out := r.Apply(record)
	if out["token"] != "***" {
		t.Errorf("expected ***, got %v", out["token"])
	}
}

func TestApplyMissingFieldIgnored(t *testing.T) {
	r, err := redact.New([]string{"secret"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{"msg": "hello"}
	out := r.Apply(record)
	if _, ok := out["secret"]; ok {
		t.Error("expected missing field to remain absent")
	}
	if out["msg"] != "hello" {
		t.Errorf("expected hello, got %v", out["msg"])
	}
}

func TestApplyBytesValid(t *testing.T) {
	r, err := redact.New([]string{"password"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := []byte(`{"user":"bob","password":"hunter2"}`)
	out, err := r.ApplyBytes(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", result["password"])
	}
	if result["user"] != "bob" {
		t.Errorf("expected bob, got %v", result["user"])
	}
}

func TestApplyBytesInvalidJSON(t *testing.T) {
	r, err := redact.New([]string{"x"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = r.ApplyBytes([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
