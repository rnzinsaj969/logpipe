package tally

import (
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

func entry(level, service, msg string) reader.LogEntry {
	return reader.LogEntry{
		Level:     level,
		Service:   service,
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestAddCountsByLevel(t *testing.T) {
	tl, err := New("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tl.Add(entry("info", "svc", "a"))
	tl.Add(entry("info", "svc", "b"))
	tl.Add(entry("error", "svc", "c"))

	snap := tl.Snapshot()
	if snap["info"] != 2 {
		t.Errorf("expected info=2, got %d", snap["info"])
	}
	if snap["error"] != 1 {
		t.Errorf("expected error=1, got %d", snap["error"])
	}
}

func TestAddCountsByService(t *testing.T) {
	tl, _ := New("service")
	tl.Add(entry("info", "alpha", "x"))
	tl.Add(entry("info", "beta", "y"))
	tl.Add(entry("warn", "alpha", "z"))

	snap := tl.Snapshot()
	if snap["alpha"] != 2 {
		t.Errorf("expected alpha=2, got %d", snap["alpha"])
	}
	if snap["beta"] != 1 {
		t.Errorf("expected beta=1, got %d", snap["beta"])
	}
}

func TestAddCountsExtraField(t *testing.T) {
	tl, _ := New("env")
	e := entry("info", "svc", "msg")
	e.Extra = map[string]interface{}{"env": "prod"}
	tl.Add(e)
	e2 := entry("info", "svc", "msg2")
	e2.Extra = map[string]interface{}{"env": "prod"}
	tl.Add(e2)

	snap := tl.Snapshot()
	if snap["prod"] != 2 {
		t.Errorf("expected prod=2, got %d", snap["prod"])
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	tl, _ := New("level")
	tl.Add(entry("info", "s", "m"))
	snap := tl.Snapshot()
	snap["info"] = 999
	snap2 := tl.Snapshot()
	if snap2["info"] != 1 {
		t.Errorf("snapshot mutation affected internal state")
	}
}

func TestResetClearsCounts(t *testing.T) {
	tl, _ := New("level")
	tl.Add(entry("info", "s", "m"))
	tl.Reset()
	snap := tl.Snapshot()
	if len(snap) != 0 {
		t.Errorf("expected empty snapshot after reset, got %v", snap)
	}
}

func TestConcurrentAdd(t *testing.T) {
	tl, _ := New("level")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Add(entry("info", "svc", "msg"))
		}()
	}
	wg.Wait()
	snap := tl.Snapshot()
	if snap["info"] != 100 {
		t.Errorf("expected 100, got %d", snap["info"])
	}
}
