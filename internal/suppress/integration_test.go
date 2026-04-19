package suppress_test

import (
	"testing"
	"time"

	"logpipe/internal/reader"
	"logpipe/internal/suppress"
)

func makeEntry(svc, msg, level string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: msg, Level: level}
}

func TestSuppressFilterSlice(t *testing.T) {
	now := time.Now()
	s, _ := suppress.New(suppress.Options{Cooldown: time.Minute})
	s.(*suppress.Suppressor) // ensure concrete type accessible via interface if needed

	entries := []reader.LogEntry{
		makeEntry("api", "started", "info"),
		makeEntry("api", "started", "info"), // duplicate
		makeEntry("api", "stopped", "info"),
		makeEntry("db", "started", "info"),
	}
	_ = now

	var passed []reader.LogEntry
	for _, e := range entries {
		if s.Apply(e) {
			passed = append(passed, e)
		}
	}
	if len(passed) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(passed))
	}
}

func TestSuppressMultipleServices(t *testing.T) {
	s, _ := suppress.New(suppress.Options{Cooldown: time.Minute})
	services := []string{"a", "b", "c"}
	for _, svc := range services {
		if !s.Apply(makeEntry(svc, "ping", "debug")) {
			t.Fatalf("first entry for %s should pass", svc)
		}
		if s.Apply(makeEntry(svc, "ping", "debug")) {
			t.Fatalf("second entry for %s should be suppressed", svc)
		}
	}
}
