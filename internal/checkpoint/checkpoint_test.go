package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNewCreatesEmptyStore(t *testing.T) {
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, ok := s.Get("http://example.com")
	if ok {
		t.Fatal("expected no entry for unknown URL")
	}
}

func TestSetAndGet(t *testing.T) {
	path := tempPath(t)
	s, _ := checkpoint.New(path)
	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	if err := s.Set("http://host/logs", 42, ts); err != nil {
		t.Fatalf("Set: %v", err)
	}

	st, ok := s.Get("http://host/logs")
	if !ok {
		t.Fatal("expected entry after Set")
	}
	if st.Offset != 42 {
		t.Errorf("offset: got %d, want 42", st.Offset)
	}
	if !st.Timestamp.Equal(ts) {
		t.Errorf("timestamp: got %v, want %v", st.Timestamp, ts)
	}
}

func TestPersistsAcrossReopen(t *testing.T) {
	path := tempPath(t)
	s1, _ := checkpoint.New(path)
	ts := time.Now().UTC().Truncate(time.Second)
	s1.Set("http://host/logs", 99, ts)

	s2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	st, ok := s2.Get("http://host/logs")
	if !ok {
		t.Fatal("expected persisted entry after reopen")
	}
	if st.Offset != 99 {
		t.Errorf("offset: got %d, want 99", st.Offset)
	}
}

func TestDelete(t *testing.T) {
	path := tempPath(t)
	s, _ := checkpoint.New(path)
	s.Set("http://host/logs", 10, time.Now())

	if err := s.Delete("http://host/logs"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get("http://host/logs")
	if ok {
		t.Fatal("expected entry to be removed after Delete")
	}
}

func TestNewMissingFileIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	_, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestNewCorruptFileReturnsError(t *testing.T) {
	path := tempPath(t)
	os.WriteFile(path, []byte("not json{"), 0o600)
	_, err := checkpoint.New(path)
	if err == nil {
		t.Fatal("expected error for corrupt checkpoint file")
	}
}
