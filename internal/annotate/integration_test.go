package annotate_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/annotate"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(msg string) reader.LogEntry {
	return reader.LogEntry{
		Service:   "api",
		Level:     "error",
		Message:   msg,
		Timestamp: time.Now(),
		Extra:     map[string]any{},
	}
}

func TestAnnotateChainMultipleEntries(t *testing.T) {
	a, err := annotate.New([]annotate.Rule{
		{Pattern: `timeout`, Key: "category", Value: "network"},
		{Pattern: `auth`, Key: "category", Value: "security"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := []reader.LogEntry{
		makeEntry("connection timeout"),
		makeEntry("auth failure"),
		makeEntry("disk full"),
	}

	expected := []string{"network", "security", ""}
	for i, e := range entries {
		out := a.Apply(e)
		got, _ := out.Extra["category"].(string)
		if got != expected[i] {
			t.Errorf("entry %d: expected category=%q, got %q", i, expected[i], got)
		}
	}
}

func TestAnnotateDoesNotMutateInput(t *testing.T) {
	a, err := annotate.New([]annotate.Rule{
		{Pattern: `.*`, Key: "matched", Value: "yes"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := makeEntry("any message")
	origExtra := len(e.Extra)
	a.Apply(e)

	if len(e.Extra) != origExtra {
		t.Error("input entry Extra was mutated")
	}
}
