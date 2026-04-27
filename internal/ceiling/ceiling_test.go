package ceiling

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "msg", Level: "info"}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAllowWithinLimit(t *testing.T) {
	c, err := New(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	c, _ := New(2, time.Second)
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny after limit reached")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	c, _ := New(1, time.Second)
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

func TestAllowResetsAfterWindow(t *testing.T) {
	c, _ := New(1, 50*time.Millisecond)
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !c.Allow(entry("svc")) {
		t.Fatal("expected allow after window reset")
	}
}
