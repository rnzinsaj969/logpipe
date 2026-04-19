package debounce

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func fixedClock(t time.Time) clock { return func() time.Time { return t } }

func entry(svc, msg string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: msg}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestFirstOccurrenceAllowed(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Second, fixedClock(now))
	if !d.Allow(entry("svc", "hello")) {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestDuplicateWithinWindowSuppressed(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Second, fixedClock(now))
	d.Allow(entry("svc", "hello"))
	if d.Allow(entry("svc", "hello")) {
		t.Fatal("expected duplicate within window to be suppressed")
	}
}

func TestDuplicateAfterWindowAllowed(t *testing.T) {
	base := time.Now()
	d, _ := newWithClock(time.Second, fixedClock(base))
	d.Allow(entry("svc", "hello"))

	d.now = fixedClock(base.Add(2 * time.Second))
	if !d.Allow(entry("svc", "hello")) {
		t.Fatal("expected entry to be allowed after window expires")
	}
}

func TestDifferentServicesAreIndependent(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Second, fixedClock(now))
	d.Allow(entry("svc-a", "hello"))
	if !d.Allow(entry("svc-b", "hello")) {
		t.Fatal("expected different service to be allowed")
	}
}

func TestEvictRemovesExpiredEntries(t *testing.T) {
	base := time.Now()
	d, _ := newWithClock(time.Second, fixedClock(base))
	d.Allow(entry("svc", "hello"))

	d.now = fixedClock(base.Add(2 * time.Second))
	d.Evict()

	if len(d.seen) != 0 {
		t.Fatalf("expected seen map to be empty after eviction, got %d", len(d.seen))
	}
}

func TestEvictRetainsFreshEntries(t *testing.T) {
	base := time.Now()
	d, _ := newWithClock(500*time.Millisecond, fixedClock(base))
	d.Allow(entry("svc", "hello"))

	d.now = fixedClock(base.Add(100 * time.Millisecond))
	d.Evict()

	if len(d.seen) != 1 {
		t.Fatalf("expected 1 entry retained, got %d", len(d.seen))
	}
}
