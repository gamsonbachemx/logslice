package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/levelfilter"
)

func TestParseLevelKnown(t *testing.T) {
	cases := []struct {
		input string
		want  levelfilter.Level
	}{
		{"debug", levelfilter.LevelDebug},
		{"INFO", levelfilter.LevelInfo},
		{"Warn", levelfilter.LevelWarn},
		{"warning", levelfilter.LevelWarn},
		{"error", levelfilter.LevelError},
		{"err", levelfilter.LevelError},
		{"fatal", levelfilter.LevelFatal},
		{"crit", levelfilter.LevelFatal},
		{"critical", levelfilter.LevelFatal},
	}
	for _, c := range cases {
		got := levelfilter.ParseLevel(c.input)
		if got != c.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestParseLevelUnknown(t *testing.T) {
	if got := levelfilter.ParseLevel("trace"); got != levelfilter.LevelUnknown {
		t.Errorf("expected LevelUnknown, got %v", got)
	}
}

func TestAllowPassesHigherLevel(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn, "level")
	line := []byte(`{"level":"error","msg":"boom"}`)
	if !f.Allow(line) {
		t.Error("expected error level to pass warn filter")
	}
}

func TestAllowBlocksLowerLevel(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn, "level")
	line := []byte(`{"level":"debug","msg":"verbose"}`)
	if f.Allow(line) {
		t.Error("expected debug level to be blocked by warn filter")
	}
}

func TestAllowPassesEqualLevel(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	line := []byte(`{"level":"error","msg":"exact"}`)
	if !f.Allow(line) {
		t.Error("expected equal level to pass")
	}
}

func TestAllowPassesInvalidJSON(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	if !f.Allow([]byte(`not-json`)) {
		t.Error("expected invalid JSON to pass through")
	}
}

func TestAllowPassesMissingField(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError, "level")
	if !f.Allow([]byte(`{"msg":"no level field"}`)) {
		t.Error("expected missing level field to pass through")
	}
}

func TestAllowCustomLevelKey(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn, "severity")
	pass := []byte(`{"severity":"fatal","msg":"crash"}`)
	if !f.Allow(pass) {
		t.Error("expected fatal to pass warn filter on custom key")
	}
	block := []byte(`{"severity":"info","msg":"routine"}`)
	if f.Allow(block) {
		t.Error("expected info to be blocked by warn filter on custom key")
	}
}

func TestNewDefaultLevelKey(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelInfo, "")
	line := []byte(`{"level":"debug","msg":"low"}`)
	if f.Allow(line) {
		t.Error("expected debug to be blocked when default key used")
	}
}
