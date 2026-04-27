package ceiling_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/ceiling"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "hello", Level: "warn"}
}

func TestConcurrentAllowIsSafe(t *testing.T) {
	c, err := ceiling.New(50, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var wg sync.WaitGroup
	var allowed int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c.Allow(makeEntry("svc")) {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}
	wg.Wait()

	if allowed > 50 {
		t.Fatalf("expected at most 50 allowed, got %d", allowed)
	}
}

func TestCeilingEnforcesLimitAcrossMultipleServices(t *testing.T) {
	c, _ := ceiling.New(2, time.Second)
	services := []string{"alpha", "beta", "gamma"}
	for _, svc := range services {
		if !c.Allow(makeEntry(svc)) {
			t.Fatalf("first call should be allowed for %s", svc)
		}
		if !c.Allow(makeEntry(svc)) {
			t.Fatalf("second call should be allowed for %s", svc)
		}
		if c.Allow(makeEntry(svc)) {
			t.Fatalf("third call should be denied for %s", svc)
		}
	}
}
