package fetch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/fetch"
)

func TestStreamSuccess(t *testing.T) {
	payload := `{"level":"info","msg":"started"}
{"level":"warn","msg":"slow query"}
{"level":"error","msg":"timeout"}
`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payload)
	}))
	defer ts.Close()

	client := fetch.NewClient(5 * time.Second)
	lines, errCh := client.Stream(fetch.Config{URL: ts.URL})

	var got []string
	for line := range lines {
		got = append(got, line)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := strings.Split(strings.TrimSpace(payload), "\n")
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: got %q, want %q", i, got[i], w)
		}
	}
}

func TestStreamNonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := fetch.NewClient(0)
	lines, errCh := client.Stream(fetch.Config{URL: ts.URL})

	for range lines {
	}
	if err := <-errCh; err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestStreamInvalidURL(t *testing.T) {
	client := fetch.NewClient(0)
	lines, errCh := client.Stream(fetch.Config{URL: "://bad-url"})

	for range lines {
	}
	if err := <-errCh; err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestNewClientDefaultTimeout(t *testing.T) {
	client := fetch.NewClient(0)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
