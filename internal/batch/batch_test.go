package batch

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(msg string) reader.LogEntry {
	return reader.LogEntry{Service: "svc", Level: "info", Message: msg}
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAddBelowMaxSizeNotReady(t *testing.T) {
	b := New(Options{MaxSize: 3, MaxAge: time.Minute})
	ready := b.Add(entry("a"))
	if ready {
		t.Fatal("expected not ready after first entry")
	}
	if b.Len() != 1 {
		t.Fatalf("expected len 1, got %d", b.Len())
	}
}

func TestAddAtMaxSizeIsReady(t *testing.T) {
	b := New(Options{MaxSize: 2, MaxAge: time.Minute})
	b.Add(entry("a"))
	ready := b.Add(entry("b"))
	if !ready {
		t.Fatal("expected ready at max size")
	}
}

func TestFlushReturnsEntriesAndResets(t *testing.T) {
	b := New(Options{MaxSize: 10, MaxAge: time.Minute})
	b.Add(entry("x"))
	b.Add(entry("y"))
	out, err := b.Flush()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if b.Len() != 0 {
		t.Fatal("expected empty batch after flush")
	}
}

func TestFlushEmptyReturnsError(t *testing.T) {
	b := New(DefaultOptions())
	_, err := b.Flush()
	if err != ErrEmptyBatch {
		t.Fatalf("expected ErrEmptyBatch, got %v", err)
	}
}

func TestReadyByAge(t *testing.T) {
	now := time.Now()
	b := New(Options{MaxSize: 100, MaxAge: 10 * time.Millisecond})
	b.clock = fixedClock(now)
	b.Add(entry("a"))

	// still within window
	if b.Ready() {
		t.Fatal("should not be ready before max age")
	}

	// advance clock past max age
	b.clock = fixedClock(now.Add(20 * time.Millisecond))
	if !b.Ready() {
		t.Fatal("should be ready after max age")
	}
}

func TestFlushIsolatesOutput(t *testing.T) {
	b := New(DefaultOptions())
	b.Add(entry("first"))
	out, _ := b.Flush()
	b.Add(entry("second"))
	if out[0].Message != "first" {
		t.Fatalf("output mutated after flush: got %q", out[0].Message)
	}
}

func TestDefaultOptionsApplied(t *testing.T) {
	b := New(Options{})
	if b.opts.MaxSize != 100 {
		t.Fatalf("expected default MaxSize 100, got %d", b.opts.MaxSize)
	}
	if b.opts.MaxAge != 5*time.Second {
		t.Fatalf("expected default MaxAge 5s, got %v", b.opts.MaxAge)
	}
}
