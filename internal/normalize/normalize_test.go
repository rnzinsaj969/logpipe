package normalize_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/normalize"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Service:   "api",
		Message:   "  hello world  ",
	}
}

func TestLowercasesLevel(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	out := n.Apply(base())
	if out.Level != "info" {
		t.Fatalf("expected 'info', got %q", out.Level)
	}
}

func TestTrimSpaceMessage(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	out := n.Apply(base())
	if out.Message != "hello world" {
		t.Fatalf("expected trimmed message, got %q", out.Message)
	}
}

func TestDefaultLevelApplied(t *testing.T) {
	opts := normalize.DefaultOptions()
	opts.DefaultLevel = "warn"
	n := normalize.New(opts)
	e := base()
	e.Level = ""
	out := n.Apply(e)
	if out.Level != "warn" {
		t.Fatalf("expected 'warn', got %q", out.Level)
	}
}

func TestDefaultServiceApplied(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	e := base()
	e.Service = ""
	out := n.Apply(e)
	if out.Service != "unknown" {
		t.Fatalf("expected 'unknown', got %q", out.Service)
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	e := base()
	origLevel := e.Level
	_ = n.Apply(e)
	if e.Level != origLevel {
		t.Fatal("original entry was mutated")
	}
}

func TestNoLowercaseWhenDisabled(t *testing.T) {
	opts := normalize.DefaultOptions()
	opts.LowercaseLevel = false
	n := normalize.New(opts)
	out := n.Apply(base())
	if out.Level != "INFO" {
		t.Fatalf("expected 'INFO', got %q", out.Level)
	}
}
