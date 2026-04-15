package merge_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/merge"
	"github.com/logpipe/logpipe/internal/reader"
)

func base(msg string) reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   msg,
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestMergeNoEntriesReturnsError(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	_, err := m.Apply(nil)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestMergeSingleEntryPassthrough(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	e := base("hello")
	e.Extra = map[string]any{"k": "v"}
	out, err := m.Apply([]reader.LogEntry{e})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Message != "hello" || out.Extra["k"] != "v" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestMergePreferFirstKeepsEarlyValue(t *testing.T) {
	m, _ := merge.New(merge.Options{PreferFirst: true})
	a := base("first")
	a.Extra = map[string]any{"key": "from-a"}
	b := base("second")
	b.Extra = map[string]any{"key": "from-b"}

	out, err := m.Apply([]reader.LogEntry{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Extra["key"] != "from-a" {
		t.Fatalf("expected from-a, got %v", out.Extra["key"])
	}
	if out.Message != "first" {
		t.Fatalf("base fields should come from first entry, got %q", out.Message)
	}
}

func TestMergePreferLastOverwritesEarlyValue(t *testing.T) {
	m, _ := merge.New(merge.Options{PreferFirst: false})
	a := base("first")
	a.Extra = map[string]any{"key": "from-a"}
	b := base("second")
	b.Extra = map[string]any{"key": "from-b"}

	out, err := m.Apply([]reader.LogEntry{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Extra["key"] != "from-b" {
		t.Fatalf("expected from-b, got %v", out.Extra["key"])
	}
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	a := base("a")
	a.Extra = map[string]any{"x": 1}
	b := base("b")
	b.Extra = map[string]any{"y": 2}

	out, _ := m.Apply([]reader.LogEntry{a, b})
	out.Extra["x"] = 99

	if a.Extra["x"] != 1 {
		t.Fatal("input a was mutated")
	}
}

func TestMergeNoExtraProducesNilMap(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	out, err := m.Apply([]reader.LogEntry{base("a"), base("b")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Extra != nil {
		t.Fatalf("expected nil Extra, got %v", out.Extra)
	}
}
