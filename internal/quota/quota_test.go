package quota

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowWithinLimit(t *testing.T) {
	q, _ := New(Options{MaxEntries: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if err := q.Allow("svc"); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i, err)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	q, _ := New(Options{MaxEntries: 2, Window: time.Minute})
	q.Allow("svc")
	q.Allow("svc")
	if err := q.Allow("svc"); err != ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	now := time.Now()
	q, _ := New(Options{MaxEntries: 1, Window: time.Second})
	q.clock = fixedClock(now)
	q.Allow("svc")
	if err := q.Allow("svc"); err != ErrQuotaExceeded {
		t.Fatal("expected exceeded before reset")
	}
	q.clock = fixedClock(now.Add(2 * time.Second))
	if err := q.Allow("svc"); err != nil {
		t.Fatalf("expected nil after window reset, got %v", err)
	}
}

func TestAllowIndependentServices(t *testing.T) {
	q, _ := New(Options{MaxEntries: 1, Window: time.Minute})
	q.Allow("a")
	if err := q.Allow("b"); err != nil {
		t.Fatalf("service b should not be limited: %v", err)
	}
}

func TestResetClearsService(t *testing.T) {
	q, _ := New(Options{MaxEntries: 1, Window: time.Minute})
	q.Allow("svc")
	q.Reset("svc")
	if err := q.Allow("svc"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestSnapshotReflectsCounts(t *testing.T) {
	q, _ := New(Options{MaxEntries: 5, Window: time.Minute})
	q.Allow("x")
	q.Allow("x")
	q.Allow("y")
	snap := q.Snapshot()
	if snap["x"] != 2 {
		t.Errorf("expected 2 for x, got %d", snap["x"])
	}
	if snap["y"] != 1 {
		t.Errorf("expected 1 for y, got %d", snap["y"])
	}
}

func TestInvalidOptionsReturnError(t *testing.T) {
	if _, err := New(Options{MaxEntries: 0, Window: time.Minute}); err == nil {
		t.Fatal("expected error for zero MaxEntries")
	}
	if _, err := New(Options{MaxEntries: 1, Window: 0}); err == nil {
		t.Fatal("expected error for zero Window")
	}
}
