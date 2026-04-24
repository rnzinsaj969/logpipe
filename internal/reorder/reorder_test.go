package reorder

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(ts time.Time, msg string) reader.LogEntry {
	return reader.LogEntry{Timestamp: ts, Message: msg, Service: "svc", Level: "info"}
}

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNewInvalidWindowSizeReturnsError(t *testing.T) {
	_, err := New(Options{WindowSize: 0, MaxAge: time.Second})
	if err == nil {
		t.Fatal("expected error for zero WindowSize")
	}
}

func TestNewInvalidMaxAgeReturnsError(t *testing.T) {
	_, err := New(Options{WindowSize: 3, MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
}

func TestAddBelowWindowSizeReturnsNil(t *testing.T) {
	r, _ := New(Options{WindowSize: 3, MaxAge: time.Second})
	out := r.Add(entry(base, "a"))
	if out != nil {
		t.Fatalf("expected nil before window full, got %v", out)
	}
}

func TestAddAtWindowSizeEmitsSorted(t *testing.T) {
	r, _ := New(Options{WindowSize: 3, MaxAge: time.Second})
	r.Add(entry(base.Add(2*time.Second), "c"))
	r.Add(entry(base.Add(1*time.Second), "b"))
	out := r.Add(entry(base, "a"))
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Message != "a" || out[1].Message != "b" || out[2].Message != "c" {
		t.Fatalf("unexpected order: %v", out)
	}
}

func TestFlushEmitsAllSorted(t *testing.T) {
	r, _ := New(Options{WindowSize: 10, MaxAge: time.Second})
	r.Add(entry(base.Add(3*time.Second), "d"))
	r.Add(entry(base, "a"))
	out := r.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Message != "a" {
		t.Fatalf("expected 'a' first, got %s", out[0].Message)
	}
}

func TestFlushResetsBuffer(t *testing.T) {
	r, _ := New(Options{WindowSize: 10, MaxAge: time.Second})
	r.Add(entry(base, "a"))
	r.Flush()
	out := r.Flush()
	if len(out) != 0 {
		t.Fatalf("expected empty after flush, got %d", len(out))
	}
}

func TestDrainEmitsExpiredEntries(t *testing.T) {
	now := base.Add(10 * time.Second)
	r, _ := newWithClock(Options{WindowSize: 10, MaxAge: 2 * time.Second}, func() time.Time { return now })
	r.Add(entry(base, "old"))
	r.Add(entry(now, "new"))
	out := r.Drain()
	if len(out) != 1 || out[0].Message != "old" {
		t.Fatalf("expected only 'old' drained, got %v", out)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.WindowSize <= 0 {
		t.Fatal("expected positive WindowSize")
	}
	if opts.MaxAge <= 0 {
		t.Fatal("expected positive MaxAge")
	}
}
