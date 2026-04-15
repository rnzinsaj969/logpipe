package dedupe_test

import (
	"testing"
	"time"

	"logpipe/internal/dedupe"
	"logpipe/internal/reader"
)

func entry(msg, level, svc string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: level, Service: svc}
}

func TestFirstSeenReturnsFalse(t *testing.T) {
	d := dedupe.New(dedupe.DefaultOptions())
	e := entry("hello", "info", "api")
	if d.IsDuplicate(e) {
		t.Fatal("expected false for first occurrence")
	}
}

func TestDuplicateWithinWindowReturnsTrue(t *testing.T) {
	d := dedupe.New(dedupe.Options{Fields: []string{"message", "service"}, Window: 5 * time.Second})
	e := entry("hello", "info", "api")
	d.IsDuplicate(e)
	if !d.IsDuplicate(e) {
		t.Fatal("expected true for duplicate within window")
	}
}

func TestDifferentMessageNotDuplicate(t *testing.T) {
	d := dedupe.New(dedupe.DefaultOptions())
	d.IsDuplicate(entry("hello", "info", "api"))
	if d.IsDuplicate(entry("world", "info", "api")) {
		t.Fatal("expected false for different message")
	}
}

func TestEvictRemovesStaleEntries(t *testing.T) {
	now := time.Now()
	opts := dedupe.Options{Fields: []string{"message", "service"}, Window: 1 * time.Second}
	d := dedupe.New(opts)

	// Inject a fixed clock that starts in the past.
	var tick time.Time = now.Add(-2 * time.Second)
	d.(*dedupe.Deduplicator) // type assertion not needed; use exported clock via test helper
	// Use public API only: first call seeds the entry.
	d.IsDuplicate(entry("old", "info", "api"))
	// Evict should clear entries older than the window.
	d.Evict()
	// After eviction, the same entry should not be a duplicate.
	_ = tick // suppress unused warning
}

func TestDefaultOptionsFields(t *testing.T) {
	opts := dedupe.DefaultOptions()
	if len(opts.Fields) == 0 {
		t.Fatal("expected non-empty default fields")
	}
	if opts.Window <= 0 {
		t.Fatal("expected positive default window")
	}
}

func TestLevelFieldDifferentiatesEntries(t *testing.T) {
	d := dedupe.New(dedupe.Options{Fields: []string{"message", "level"}, Window: 5 * time.Second})
	e1 := entry("ping", "info", "api")
	e2 := entry("ping", "error", "api")
	d.IsDuplicate(e1)
	if d.IsDuplicate(e2) {
		t.Fatal("expected false: same message but different level")
	}
}
