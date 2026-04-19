package fork_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/fork"
	"github.com/logpipe/logpipe/internal/reader"
)

func base(level string) reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   "svc",
		Message:   "msg",
	}
}

func TestNewNilPredicateReturnsError(t *testing.T) {
	_, err := fork.New(nil, func(reader.LogEntry) error { return nil }, func(reader.LogEntry) error { return nil })
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestNewNilLeftReturnsError(t *testing.T) {
	_, err := fork.New(func(reader.LogEntry) bool { return true }, nil, func(reader.LogEntry) error { return nil })
	if err == nil {
		t.Fatal("expected error for nil left sink")
	}
}

func TestNewNilRightReturnsError(t *testing.T) {
	_, err := fork.New(func(reader.LogEntry) bool { return true }, func(reader.LogEntry) error { return nil }, nil)
	if err == nil {
		t.Fatal("expected error for nil right sink")
	}
}

func TestApplyRoutesToLeft(t *testing.T) {
	var got string
	f, _ := fork.New(
		func(e reader.LogEntry) bool { return e.Level == "error" },
		func(e reader.LogEntry) error { got = "left"; return nil },
		func(e reader.LogEntry) error { got = "right"; return nil },
	)
	_ = f.Apply(base("error"))
	if got != "left" {
		t.Fatalf("expected left, got %s", got)
	}
}

func TestApplyRoutesToRight(t *testing.T) {
	var got string
	f, _ := fork.New(
		func(e reader.LogEntry) bool { return e.Level == "error" },
		func(e reader.LogEntry) error { got = "left"; return nil },
		func(e reader.LogEntry) error { got = "right"; return nil },
	)
	_ = f.Apply(base("info"))
	if got != "right" {
		t.Fatalf("expected right, got %s", got)
	}
}

func TestApplyPropagatesSinkError(t *testing.T) {
	want := errors.New("sink error")
	f, _ := fork.New(
		func(e reader.LogEntry) bool { return true },
		func(e reader.LogEntry) error { return want },
		func(e reader.LogEntry) error { return nil },
	)
	if err := f.Apply(base("info")); !errors.Is(err, want) {
		t.Fatalf("expected sink error, got %v", err)
	}
}
