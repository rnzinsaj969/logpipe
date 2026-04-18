package group_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/group"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(level, service string, extra map[string]any) reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   service,
		Message:   "msg",
		Extra:     extra,
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := group.New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestGroupByLevel(t *testing.T) {
	g, _ := group.New("level")
	g.Add(entry("info", "svc", nil))
	g.Add(entry("error", "svc", nil))
	g.Add(entry("info", "svc", nil))

	snap := g.Snapshot()
	if len(snap["info"]) != 2 {
		t.Fatalf("expected 2 info entries, got %d", len(snap["info"]))
	}
	if len(snap["error"]) != 1 {
		t.Fatalf("expected 1 error entry, got %d", len(snap["error"]))
	}
}

func TestGroupByService(t *testing.T) {
	g, _ := group.New("service")
	g.Add(entry("info", "alpha", nil))
	g.Add(entry("info", "beta", nil))
	g.Add(entry("info", "alpha", nil))

	snap := g.Snapshot()
	if len(snap["alpha"]) != 2 {
		t.Fatalf("expected 2 alpha entries, got %d", len(snap["alpha"]))
	}
}

func TestGroupByExtraField(t *testing.T) {
	g, _ := group.New("region")
	g.Add(entry("info", "svc", map[string]any{"region": "us-east"}))
	g.Add(entry("info", "svc", map[string]any{"region": "eu-west"}))
	g.Add(entry("info", "svc", map[string]any{"region": "us-east"}))

	snap := g.Snapshot()
	if len(snap["us-east"]) != 2 {
		t.Fatalf("expected 2 us-east entries, got %d", len(snap["us-east"]))
	}
}

func TestSnapshotResetsState(t *testing.T) {
	g, _ := group.New("level")
	g.Add(entry("info", "svc", nil))
	g.Snapshot()
	snap := g.Snapshot()
	if len(snap) != 0 {
		t.Fatal("expected empty snapshot after reset")
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	g, _ := group.New("level")
	g.Add(entry("info", "svc", nil))
	snap := g.Snapshot()
	snap["info"] = nil
	g.Add(entry("info", "svc", nil))
	snap2 := g.Snapshot()
	if len(snap2["info"]) != 1 {
		t.Fatal("snapshot mutation affected internal state")
	}
}
