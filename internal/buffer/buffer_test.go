package buffer_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/buffer"
	"github.com/yourorg/logpipe/internal/filter"
)

func entry(msg string) filter.LogEntry {
	return filter.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestPushAndDrain(t *testing.T) {
	b := buffer.New(3)
	b.Push(entry("a"))
	b.Push(entry("b"))
	b.Push(entry("c"))

	out := b.Drain()
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Message != "a" || out[1].Message != "b" || out[2].Message != "c" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestDrainResetsBuffer(t *testing.T) {
	b := buffer.New(4)
	b.Push(entry("x"))
	b.Drain()
	if b.Len() != 0 {
		t.Errorf("expected empty buffer after drain, got %d", b.Len())
	}
}

func TestRingOverwrite(t *testing.T) {
	b := buffer.New(2)
	b.Push(entry("first"))
	b.Push(entry("second"))
	b.Push(entry("third")) // overwrites "first"

	out := b.Drain()
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Message != "second" || out[1].Message != "third" {
		t.Errorf("expected second,third got %s,%s", out[0].Message, out[1].Message)
	}
}

func TestLen(t *testing.T) {
	b := buffer.New(10)
	if b.Len() != 0 {
		t.Errorf("expected 0, got %d", b.Len())
	}
	b.Push(entry("one"))
	b.Push(entry("two"))
	if b.Len() != 2 {
		t.Errorf("expected 2, got %d", b.Len())
	}
}

func TestPanicOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for capacity 0")
		}
	}()
	buffer.New(0)
}
