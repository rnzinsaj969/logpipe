package cap_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/cap"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Level: "info", Message: "hello"}
}

func TestAllowWithinLimit(t *testing.T) {
	c, _ := cap.New(cap.Options{Max: 3, Window: time.Second})
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	c, _ := cap.New(cap.Options{Max: 2, Window: time.Second})
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny on third call")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	now := time.Now()
	c, _ := cap.New(cap.Options{Max: 1, Window: 50 * time.Millisecond})
	// patch internal clock via a subtest helper — use real sleep instead
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny before window expires")
	}
	time.Sleep(60 * time.Millisecond)
	_ = now
	if !c.Allow(entry("svc")) {
		t.Fatal("expected allow after window reset")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	c, _ := cap.New(cap.Options{Max: 1, Window: time.Second})
	if !c.Allow(entry("a")) {
		t.Fatal("expected allow for service a")
	}
	if !c.Allow(entry("b")) {
		t.Fatal("expected allow for service b")
	}
	if c.Allow(entry("a")) {
		t.Fatal("expected deny for service a on second call")
	}
}

func TestInvalidOptionsReturnError(t *testing.T) {
	if _, err := cap.New(cap.Options{Max: 0, Window: time.Second}); err == nil {
		t.Fatal("expected error for Max=0")
	}
	if _, err := cap.New(cap.Options{Max: 1, Window: 0}); err == nil {
		t.Fatal("expected error for Window=0")
	}
}

func TestConcurrentAllow(t *testing.T) {
	c, _ := cap.New(cap.Options{Max: 50, Window: time.Second})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); c.Allow(entry("svc")) }()
	}
	wg.Wait()
}
