package health_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/health"
)

func TestMonitorConcurrentAccess(t *testing.T) {
	m := health.New()
	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			name := "svc"
			if idx%2 == 0 {
				m.RecordSuccess(name)
			} else {
				m.RecordError(name)
			}
		}(i)
	}

	wg.Wait()

	snap := m.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 source, got %d", len(snap))
	}
	if snap[0].ErrorCount < 0 {
		t.Error("error count must not be negative")
	}
}

func TestMonitorMultipleSources(t *testing.T) {
	m := health.New()
	services := []string{"auth", "billing", "gateway"}

	for _, svc := range services {
		m.RecordSuccess(svc)
	}
	m.RecordError("billing")

	snap := m.Snapshot()
	if len(snap) != len(services) {
		t.Fatalf("expected %d sources, got %d", len(services), len(snap))
	}

	statuses := make(map[string]health.Status)
	for _, s := range snap {
		statuses[s.Name] = s.Status
	}

	if statuses["auth"] != health.StatusOK {
		t.Errorf("auth: expected OK, got %s", statuses["auth"])
	}
	if statuses["billing"] != health.StatusDegraded {
		t.Errorf("billing: expected degraded, got %s", statuses["billing"])
	}
	if statuses["gateway"] != health.StatusOK {
		t.Errorf("gateway: expected OK, got %s", statuses["gateway"])
	}
}
