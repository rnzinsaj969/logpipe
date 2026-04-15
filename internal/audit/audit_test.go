package audit_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/audit"
)

func TestRecordAndSnapshot(t *testing.T) {
	l := audit.New(10)
	l.Record(audit.EventDropped, "svc-a", "rate limit", nil)
	l.Record(audit.EventRouted, "svc-b", "matched rule", map[string]string{"dest": "stdout"})

	events := l.Snapshot()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Kind != audit.EventDropped {
		t.Errorf("expected dropped, got %s", events[0].Kind)
	}
	if events[1].Meta["dest"] != "stdout" {
		t.Errorf("expected meta dest=stdout, got %v", events[1].Meta)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	l := audit.New(10)
	l.Record(audit.EventRedacted, "svc", "bearer token", nil)

	snap := l.Snapshot()
	snap[0].Service = "mutated"

	original := l.Snapshot()
	if original[0].Service == "mutated" {
		t.Error("snapshot mutation affected internal state")
	}
}

func TestCapacityEvictsOldest(t *testing.T) {
	l := audit.New(3)
	for i := 0; i < 5; i++ {
		l.Record(audit.EventTransformed, "svc", "step", nil)
	}
	if l.Len() != 3 {
		t.Errorf("expected len 3 after eviction, got %d", l.Len())
	}
}

func TestClear(t *testing.T) {
	l := audit.New(10)
	l.Record(audit.EventDropped, "svc", "filtered", nil)
	l.Clear()
	if l.Len() != 0 {
		t.Errorf("expected empty log after clear, got %d", l.Len())
	}
}

func TestConcurrentRecord(t *testing.T) {
	l := audit.New(512)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Record(audit.EventRouted, "svc", "concurrent", nil)
		}()
	}
	wg.Wait()
	if l.Len() != 100 {
		t.Errorf("expected 100 events, got %d", l.Len())
	}
}

func TestDefaultCapacity(t *testing.T) {
	l := audit.New(0) // should default to 256
	for i := 0; i < 300; i++ {
		l.Record(audit.EventDropped, "svc", "overflow", nil)
	}
	if l.Len() != 256 {
		t.Errorf("expected default capacity 256, got %d", l.Len())
	}
}
