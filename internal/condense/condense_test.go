package condense

import (
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/reader"
)

func entry(svc, msg string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: msg, Level: "info"}
}

func TestNewInvalidMinPrefixReturnsError(t *testing.T) {
	_, err := New(Options{MinPrefix: 0, MaxAge: time.Second})
	if err == nil {
		t.Fatal("expected error for MinPrefix=0")
	}
}

func TestNewInvalidMaxAgeReturnsError(t *testing.T) {
	_, err := New(Options{MinPrefix: 4, MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for MaxAge=0")
	}
}

func TestFirstEntryReturnsNil(t *testing.T) {
	c, _ := New(DefaultOptions())
	out := c.Apply(entry("svc", "starting up server"))
	if out != nil {
		t.Fatalf("expected nil for first entry, got %v", out)
	}
}

func TestMatchingPrefixAbsorbsEntry(t *testing.T) {
	c, _ := New(DefaultOptions())
	c.Apply(entry("svc", "starting up server on port 8080"))
	out := c.Apply(entry("svc", "starting up server on port 9090"))
	if out != nil {
		t.Fatalf("expected absorbed (nil), got %v", out)
	}
}

func TestShortPrefixFlushesGroup(t *testing.T) {
	c, _ := New(Options{MinPrefix: 8, MaxAge: time.Second})
	c.Apply(entry("svc", "starting up server"))
	out := c.Apply(entry("svc", "xyz"))
	if out == nil {
		t.Fatal("expected flushed entry when prefix too short")
	}
	if out.Service != "svc" {
		t.Fatalf("unexpected service: %s", out.Service)
	}
}

func TestCondensedMessageIncludesCount(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := New(Options{MinPrefix: 4, MaxAge: 5 * time.Second})
	c.clock = func() time.Time { return now }

	c.Apply(entry("svc", "processing request id=1"))
	c.Apply(entry("svc", "processing request id=2"))
	c.Apply(entry("svc", "processing request id=3"))

	// advance past MaxAge to flush
	c.clock = func() time.Time { return now.Add(10 * time.Second) }
	entries := c.Flush()
	if len(entries) != 1 {
		t.Fatalf("expected 1 flushed entry, got %d", len(entries))
	}
	got := entries[0].Message
	expect := "processing request id= … (3x)"
	if got != expect {
		t.Fatalf("expected %q, got %q", expect, got)
	}
}

func TestFlushReturnsNothingBeforeMaxAge(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := New(Options{MinPrefix: 4, MaxAge: 5 * time.Second})
	c.clock = func() time.Time { return now }
	c.Apply(entry("svc", "hello world"))

	out := c.Flush()
	if len(out) != 0 {
		t.Fatalf("expected empty flush before MaxAge, got %d entries", len(out))
	}
}

func TestDifferentServicesAreIndependent(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := New(Options{MinPrefix: 4, MaxAge: 5 * time.Second})
	c.clock = func() time.Time { return now }

	c.Apply(entry("alpha", "loading config file"))
	c.Apply(entry("beta", "loading config module"))

	c.clock = func() time.Time { return now.Add(10 * time.Second) }
	out := c.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2 independent groups flushed, got %d", len(out))
	}
}

func TestSingleEntryGroupPreservesOriginalMessage(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := New(Options{MinPrefix: 4, MaxAge: time.Second})
	c.clock = func() time.Time { return now }
	c.Apply(entry("svc", "unique message"))

	c.clock = func() time.Time { return now.Add(2 * time.Second) }
	out := c.Flush()
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Message != "unique message" {
		t.Fatalf("expected original message, got %q", out[0].Message)
	}
}
