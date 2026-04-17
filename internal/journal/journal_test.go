package journal_test

import (
	"testing"
	"time"

	"github.com/logpipe/internal/journal"
	"github.com/logpipe/internal/reader"
)

func entry(msg string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: "info", Service: "svc", Timestamp: time.Now()}
}

func TestInvalidCapacityReturnsError(t *testing.T) {
	_, err := journal.New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestAppendAndLen(t *testing.T) {
	j, _ := journal.New(10)
	j.Append("a", entry("hello"))
	j.Append("b", entry("world"))
	if j.Len() != 2 {
		t.Fatalf("expected 2, got %d", j.Len())
	}
}

func TestSequenceMonotonicallyIncreases(t *testing.T) {
	j, _ := journal.New(10)
	e1 := j.Append("s", entry("a"))
	e2 := j.Append("s", entry("b"))
	if e2.Seq <= e1.Seq {
		t.Fatalf("seq not increasing: %d <= %d", e2.Seq, e1.Seq)
	}
}

func TestEvictsOldestWhenFull(t *testing.T) {
	j, _ := journal.New(3)
	j.Append("s", entry("first"))
	j.Append("s", entry("second"))
	j.Append("s", entry("third"))
	j.Append("s", entry("fourth"))
	snap := j.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3, got %d", len(snap))
	}
	if snap[0].Log.Message != "second" {
		t.Fatalf("expected 'second', got %q", snap[0].Log.Message)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	j, _ := journal.New(10)
	j.Append("s", entry("x"))
	snap := j.Snapshot()
	j.Append("s", entry("y"))
	if len(snap) != 1 {
		t.Fatal("snapshot should not reflect later appends")
	}
}

func TestClearResetsLen(t *testing.T) {
	j, _ := journal.New(10)
	j.Append("s", entry("a"))
	j.Append("s", entry("b"))
	j.Clear()
	if j.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", j.Len())
	}
}

func TestSourceTagPreserved(t *testing.T) {
	j, _ := journal.New(10)
	j.Append("my-service", entry("msg"))
	snap := j.Snapshot()
	if snap[0].Source != "my-service" {
		t.Fatalf("expected 'my-service', got %q", snap[0].Source)
	}
}
