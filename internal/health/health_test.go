package health

import (
	"testing"
	"time"
)

func TestRecordSuccessSetsOK(t *testing.T) {
	m := New()
	m.RecordSuccess("svc-a")

	snap := m.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 source, got %d", len(snap))
	}
	if snap[0].Status != StatusOK {
		t.Errorf("expected status OK, got %s", snap[0].Status)
	}
	if snap[0].LastSeen.IsZero() {
		t.Error("expected LastSeen to be set")
	}
}

func TestRecordErrorDegraded(t *testing.T) {
	m := New()
	m.RecordError("svc-b")
	m.RecordError("svc-b")

	snap := m.Snapshot()
	if snap[0].Status != StatusDegraded {
		t.Errorf("expected degraded, got %s", snap[0].Status)
	}
	if snap[0].ErrorCount != 2 {
		t.Errorf("expected error count 2, got %d", snap[0].ErrorCount)
	}
}

func TestRecordErrorDown(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordError("svc-c")
	}

	snap := m.Snapshot()
	if snap[0].Status != StatusDown {
		t.Errorf("expected down, got %s", snap[0].Status)
	}
}

func TestRecordSuccessResetsStatus(t *testing.T) {
	m := New()
	for i := 0; i < 5; i++ {
		m.RecordError("svc-d")
	}
	m.RecordSuccess("svc-d")

	snap := m.Snapshot()
	if snap[0].Status != StatusOK {
		t.Errorf("expected OK after success, got %s", snap[0].Status)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	m := New()
	m.RecordSuccess("svc-e")

	snap1 := m.Snapshot()
	snap1[0].Status = StatusDown

	snap2 := m.Snapshot()
	if snap2[0].Status != StatusOK {
		t.Error("snapshot mutation should not affect monitor state")
	}
}

func TestLastSeenUpdated(t *testing.T) {
	m := New()
	m.RecordSuccess("svc-f")
	before := m.Snapshot()[0].LastSeen

	time.Sleep(2 * time.Millisecond)
	m.RecordSuccess("svc-f")
	after := m.Snapshot()[0].LastSeen

	if !after.After(before) {
		t.Error("expected LastSeen to advance on second success")
	}
}
