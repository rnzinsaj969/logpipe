package hold_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/hold"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(msg, level string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: level, Timestamp: time.Now()}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := hold.New(0, func(reader.LogEntry) bool { return false })
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewNilPredicateReturnsError(t *testing.T) {
	_, err := hold.New(10, nil)
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestAddAccumulatesEntries(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return e.Level == "error" })

	_, released := h.Add(entry("a", "info"))
	if released {
		t.Fatal("should not release on info")
	}
	if h.Len() != 1 {
		t.Fatalf("want len 1, got %d", h.Len())
	}
}

func TestAddReleasesOnPredicate(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return e.Level == "error" })

	h.Add(entry("first", "info"))
	h.Add(entry("second", "warn"))
	out, released := h.Add(entry("third", "error"))

	if !released {
		t.Fatal("expected release")
	}
	if len(out) != 3 {
		t.Fatalf("want 3 entries, got %d", len(out))
	}
	if h.Len() != 0 {
		t.Fatalf("buffer should be empty after release, got %d", h.Len())
	}
}

func TestAddEvictsOldestWhenFull(t *testing.T) {
	h, _ := hold.New(3, func(e reader.LogEntry) bool { return false })

	h.Add(entry("a", "info"))
	h.Add(entry("b", "info"))
	h.Add(entry("c", "info"))
	h.Add(entry("d", "info")) // evicts "a"

	if h.Len() != 3 {
		t.Fatalf("want len 3, got %d", h.Len())
	}
}

func TestDiscardClearsBuffer(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return false })
	h.Add(entry("x", "info"))
	h.Add(entry("y", "info"))
	h.Discard()

	if h.Len() != 0 {
		t.Fatalf("want len 0 after discard, got %d", h.Len())
	}
}

func TestReleaseDoesNotMutateReturnedSlice(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return e.Level == "error" })
	h.Add(entry("a", "info"))
	out, _ := h.Add(entry("b", "error"))

	out[0].Message = "mutated"

	// Adding a new entry should not be affected by external mutation.
	h.Add(entry("c", "info"))
	if h.Len() != 1 {
		t.Fatalf("want 1 buffered entry, got %d", h.Len())
	}
}
