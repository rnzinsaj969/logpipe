package drop_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/drop"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(level, service, message string) reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   service,
		Message:   message,
	}
}

func TestDropByLevel(t *testing.T) {
	d, err := drop.New([]drop.Rule{{Level: "debug"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.ShouldDrop(entry("debug", "svc", "msg")) {
		t.Error("expected debug entry to be dropped")
	}
	if d.ShouldDrop(entry("info", "svc", "msg")) {
		t.Error("expected info entry to be kept")
	}
}

func TestDropByService(t *testing.T) {
	d, err := drop.New([]drop.Rule{{Service: "noisy"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.ShouldDrop(entry("info", "noisy", "hello")) {
		t.Error("expected noisy service entry to be dropped")
	}
	if d.ShouldDrop(entry("info", "quiet", "hello")) {
		t.Error("expected quiet service entry to be kept")
	}
}

func TestDropByPattern(t *testing.T) {
	d, err := drop.New([]drop.Rule{{Pattern: `healthcheck`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !entry("info", "svc", "GET /healthcheck 200")) {
		t.Error("expected healthcheck entry to be dropped")
	}
	if d.ShouldDrop(entry("info", "svc", "user logged in")) {
		t.Error("expected non-matching entry to be kept")
	}
}

func TestDropCombinedRuleAllFieldsMustMatch(t *testing.T) {
	d, err := drop.New([]drop.Rule{{Level: "warn", Service: "api", Pattern: `timeout`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// all three match
	if !d.ShouldDrop(entry("warn", "api", "connection timeout")) {
		t.Error("expected combined-match entry to be dropped")
	}
	// wrong level
	if d.ShouldDrop(entry("error", "api", "connection timeout")) {
		t.Error("expected entry with wrong level to be kept")
	}
}

func TestApplyFiltersSlice(t *testing.T) {
	d, err := drop.New([]drop.Rule{{Level: "debug"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := []reader.LogEntry{
		entry("debug", "svc", "verbose"),
		entry("info", "svc", "important"),
		entry("debug", "svc", "more noise"),
	}
	out := d.Apply(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Message != "important" {
		t.Errorf("unexpected message: %s", out[0].Message)
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := drop.New([]drop.Rule{{Pattern: `[invalid`}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}
