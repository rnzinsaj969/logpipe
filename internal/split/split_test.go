package split_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/split"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "batch",
		Timestamp: time.Unix(1000, 0),
		Extra:     map[string]any{"items": []any{"a", "b", "c"}},
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := split.New(split.Options{})
	if err == nil {
		t.Fatal("expected error for empty Field, got nil")
	}
}

func TestApplyExpandsSliceField(t *testing.T) {
	s, _ := split.New(split.Options{Field: "items"})
	out := s.Apply(base())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	for i, entry := range out {
		want := []string{"a", "b", "c"}[i]
		if entry.Extra["items"] != want {
			t.Errorf("entry %d: items = %v, want %s", i, entry.Extra["items"], want)
		}
	}
}

func TestApplyKeepOriginal(t *testing.T) {
	s, _ := split.New(split.Options{Field: "items", KeepOriginal: true})
	out := s.Apply(base())
	// 1 original + 3 split = 4
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	// first entry should be the original (items is still a slice)
	if _, ok := out[0].Extra["items"].([]any); !ok {
		t.Error("first entry should retain original slice field")
	}
}

func TestApplyMissingFieldReturnsOriginal(t *testing.T) {
	s, _ := split.New(split.Options{Field: "missing"})
	e := base()
	out := s.Apply(e)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Message != e.Message {
		t.Error("returned entry should equal original")
	}
}

func TestApplyNonSliceFieldReturnsOriginal(t *testing.T) {
	s, _ := split.New(split.Options{Field: "items"})
	e := base()
	e.Extra["items"] = "not-a-slice"
	out := s.Apply(e)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	s, _ := split.New(split.Options{Field: "items"})
	e := base()
	out := s.Apply(e)
	// mutate a cloned entry
	out[0].Extra["items"] = "mutated"
	if e.Extra["items"].([]any)[0] != "a" {
		t.Error("original entry was mutated")
	}
}
