package metrics_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/metrics"
)

func TestCounterIncrements(t *testing.T) {
	var c metrics.Counter
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounterConcurrentIncrements(t *testing.T) {
	var c metrics.Counter
	var wg sync.WaitGroup
	const goroutines = 100
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, c.Value())
	}
}

func TestMetricsSnapshot(t *testing.T) {
	m := metrics.New()
	m.EntriesRead.Inc()
	m.EntriesRead.Inc()
	m.EntriesMatched.Inc()
	m.ParseErrors.Inc()

	snap := m.Snapshot()
	if snap.EntriesRead != 2 {
		t.Errorf("EntriesRead: expected 2, got %d", snap.EntriesRead)
	}
	if snap.EntriesMatched != 1 {
		t.Errorf("EntriesMatched: expected 1, got %d", snap.EntriesMatched)
	}
	if snap.EntriesDropped != 0 {
		t.Errorf("EntriesDropped: expected 0, got %d", snap.EntriesDropped)
	}
	if snap.ParseErrors != 1 {
		t.Errorf("ParseErrors: expected 1, got %d", snap.ParseErrors)
	}
}

func TestMetricsSnapshotIsIsolated(t *testing.T) {
	m := metrics.New()
	m.EntriesRead.Inc()
	snap1 := m.Snapshot()
	m.EntriesRead.Inc()
	snap2 := m.Snapshot()
	if snap1.EntriesRead != 1 {
		t.Errorf("snap1 EntriesRead: expected 1, got %d", snap1.EntriesRead)
	}
	if snap2.EntriesRead != 2 {
		t.Errorf("snap2 EntriesRead: expected 2, got %d", snap2.EntriesRead)
	}
}
