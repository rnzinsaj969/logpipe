package suppress

import (
	"testing"
	"time"

	"logpipe/internal/reader"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func entry(svc, msg string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: msg, Level: "info"}
}

func TestFirstOccurrenceAllowed(t *testing.T) {
	s, _ := New(DefaultOptions())
	if !s.Apply(entry("svc", "hello")) {
		t.Fatal("expected first entry to pass")
	}
}

func TestDuplicateWithinCooldownSuppressed(t *testing.T) {
	now := time.Now()
	s, _ := New(DefaultOptions())
	s.clock = fixedClock(now)
	s.Apply(entry("svc", "hello"))
	if s.Apply(entry("svc", "hello")) {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestDuplicateAfterCooldownAllowed(t *testing.T) {
	now := time.Now()
	s, _ := New(Options{Cooldown: time.Second})
	s.clock = fixedClock(now)
	s.Apply(entry("svc", "hello"))
	s.clock = fixedClock(now.Add(2 * time.Second))
	if !s.Apply(entry("svc", "hello")) {
		t.Fatal("expected entry after cooldown to pass")
	}
}

func TestDifferentServiceNotSuppressed(t *testing.T) {
	now := time.Now()
	s, _ := New(DefaultOptions())
	s.clock = fixedClock(now)
	s.Apply(entry("svc-a", "hello"))
	if !s.Apply(entry("svc-b", "hello")) {
		t.Fatal("expected different service to pass")
	}
}

func TestEvictClearsStaleEntries(t *testing.T) {
	now := time.Now()
	s, _ := New(Options{Cooldown: time.Second})
	s.clock = fixedClock(now)
	s.Apply(entry("svc", "hello"))
	s.clock = fixedClock(now.Add(2 * time.Second))
	s.Evict()
	if len(s.seen) != 0 {
		t.Fatalf("expected empty seen map, got %d entries", len(s.seen))
	}
}

func TestInvalidCooldownReturnsError(t *testing.T) {
	_, err := New(Options{Cooldown: 0})
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}
