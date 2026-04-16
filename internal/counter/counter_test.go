package counter_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/counter"
)

func TestIncAndGet(t *testing.T) {
	c := counter.New()
	c.Inc("svc-a")
	c.Inc("svc-a")
	c.Inc("svc-b")
	if got := c.Get("svc-a"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	if got := c.Get("svc-b"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestAddIgnoresNonPositive(t *testing.T) {
	c := counter.New()
	c.Add("svc", 0)
	c.Add("svc", -5)
	if got := c.Get("svc"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	c := counter.New()
	c.Inc("svc")
	snap := c.Snapshot()
	snap["svc"] = 999
	if got := c.Get("svc"); got != 1 {
		t.Fatalf("snapshot mutation affected counter: got %d", got)
	}
}

func TestReset(t *testing.T) {
	c := counter.New()
	c.Inc("svc")
	c.Reset()
	if got := c.Get("svc"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestConcurrentInc(t *testing.T) {
	c := counter.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("svc")
		}()
	}
	wg.Wait()
	if got := c.Get("svc"); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}
