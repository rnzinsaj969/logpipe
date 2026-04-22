package limiter_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/internal/limiter"
	"github.com/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{
		Service:   svc,
		Message:   "test",
		Level:     "info",
		Timestamp: time.Now(),
	}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := limiter.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAcquireWithinLimit(t *testing.T) {
	l, _ := limiter.New(2)
	if !l.Acquire(entry("svc")) {
		t.Fatal("expected first acquire to succeed")
	}
	if !l.Acquire(entry("svc")) {
		t.Fatal("expected second acquire to succeed")
	}
}

func TestAcquireExceedsLimit(t *testing.T) {
	l, _ := limiter.New(1)
	if !l.Acquire(entry("svc")) {
		t.Fatal("expected first acquire to succeed")
	}
	if l.Acquire(entry("svc")) {
		t.Fatal("expected second acquire to fail")
	}
}

func TestReleaseFreesSlot(t *testing.T) {
	l, _ := limiter.New(1)
	e := entry("svc")
	l.Acquire(e)
	l.Release(e)
	if !l.Acquire(e) {
		t.Fatal("expected acquire to succeed after release")
	}
}

func TestReleaseNoop(t *testing.T) {
	l, _ := limiter.New(2)
	// Release without acquire should not panic or go negative.
	l.Release(entry("svc"))
	snap := l.Snapshot()
	if snap["svc"] != 0 {
		t.Fatalf("expected 0, got %d", snap["svc"])
	}
}

func TestIndependentServices(t *testing.T) {
	l, _ := limiter.New(1)
	if !l.Acquire(entry("a")) {
		t.Fatal("expected acquire for 'a' to succeed")
	}
	if !l.Acquire(entry("b")) {
		t.Fatal("expected acquire for 'b' to succeed independently")
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	l, _ := limiter.New(5)
	l.Acquire(entry("svc"))
	snap := l.Snapshot()
	snap["svc"] = 99
	snap2 := l.Snapshot()
	if snap2["svc"] != 1 {
		t.Fatalf("expected 1, got %d", snap2["svc"])
	}
}

func TestConcurrentAcquireRelease(t *testing.T) {
	l, _ := limiter.New(10)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := entry("svc")
			if l.Acquire(e) {
				l.Release(e)
			}
		}()
	}
	wg.Wait()
	snap := l.Snapshot()
	if snap["svc"] != 0 {
		t.Fatalf("expected 0 after all releases, got %d", snap["svc"])
	}
}
