package pivot

import (
	"testing"
	"time"

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
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestGroupByLevel(t *testing.T) {
	tbl, _ := New("level")
	tbl.Add(entry("info", "svc", nil))
	tbl.Add(entry("info", "svc", nil))
	tbl.Add(entry("error", "svc", nil))

	snap := tbl.Snapshot()
	if snap["info"] != 2 {
		t.Errorf("expected info=2, got %d", snap["info"])
	}
	if snap["error"] != 1 {
		t.Errorf("expected error=1, got %d", snap["error"])
	}
}

func TestGroupByService(t *testing.T) {
	tbl, _ := New("service")
	tbl.Add(entry("info", "api", nil))
	tbl.Add(entry("info", "worker", nil))
	tbl.Add(entry("info", "api", nil))

	snap := tbl.Snapshot()
	if snap["api"] != 2 {
		t.Errorf("expected api=2, got %d", snap["api"])
	}
	if snap["worker"] != 1 {
		t.Errorf("expected worker=1, got %d", snap["worker"])
	}
}

func TestGroupByExtraField(t *testing.T) {
	tbl, _ := New("region")
	tbl.Add(entry("info", "svc", map[string]any{"region": "us-east"}))
	tbl.Add(entry("info", "svc", map[string]any{"region": "eu-west"}))
	tbl.Add(entry("info", "svc", map[string]any{"region": "us-east"}))

	snap := tbl.Snapshot()
	if snap["us-east"] != 2 {
		t.Errorf("expected us-east=2, got %d", snap["us-east"])
	}
}

func TestSnapshotResetsState(t *testing.T) {
	tbl, _ := New("level")
	tbl.Add(entry("info", "svc", nil))
	tbl.Snapshot()
	if tbl.Len() != 0 {
		t.Errorf("expected 0 after snapshot, got %d", tbl.Len())
	}
}

func TestMissingExtraKeyFallsBackToEmpty(t *testing.T) {
	tbl, _ := New("env")
	tbl.Add(entry("info", "svc", nil))
	snap := tbl.Snapshot()
	if snap[""] != 1 {
		t.Errorf("expected empty key count=1, got %d", snap[""])
	}
}
