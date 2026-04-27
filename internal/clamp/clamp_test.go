package clamp

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "msg", Level: "info"}
}

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestAllowWithinMax(t *testing.T) {
	now := time.Now()
	c, err := newWithClock(3, time.Second, fixedClock(now))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllowExceedsMax(t *testing.T) {
	now := time.Now()
	c, _ := newWithClock(2, time.Second, fixedClock(now))
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected Allow=false when limit exceeded")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	now := time.Now()
	c, _ := newWithClock(1, time.Second, fixedClock(now))
	c.Allow(entry("svc"))
	// Advance past the window.
	c.clock = fixedClock(now.Add(2 * time.Second))
	if !c.Allow(entry("svc")) {
		t.Fatal("expected Allow=true after window reset")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	now := time.Now()
	c, _ := newWithClock(1, time.Second, fixedClock(now))
	if !c.Allow(entry("a")) {
		t.Fatal("expected Allow=true for service a")
	}
	if !c.Allow(entry("b")) {
		t.Fatal("expected Allow=true for service b")
	}
	if c.Allow(entry("a")) {
		t.Fatal("expected Allow=false for service a on second call")
	}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(1, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}
