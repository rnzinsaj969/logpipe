package clamp

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "msg", Level: "info"}
}

type fixedClock struct{ t time.Time }

func (f *fixedClock) Now() time.Time { return f.t }

func TestAllowWithinMax(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1000, 0)}
	c, err := newWithClock(3, time.Second, clk)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := entry("svc-a")
	for i := 0; i < 3; i++ {
		if !c.Allow(e) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsMax(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1000, 0)}
	c, _ := newWithClock(2, time.Second, clk)
	e := entry("svc-b")
	c.Allow(e)
	c.Allow(e)
	if c.Allow(e) {
		t.Fatal("expected entry to be dropped after exceeding max")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1000, 0)}
	c, _ := newWithClock(2, time.Second, clk)
	e := entry("svc-c")
	c.Allow(e)
	c.Allow(e)
	// advance past window
	clk.t = time.Unix(1002, 0)
	if !c.Allow(e) {
		t.Fatal("expected allow after window reset")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	clk := &fixedClock{t: time.Unix(1000, 0)}
	c, _ := newWithClock(1, time.Second, clk)
	if !c.Allow(entry("svc-x")) {
		t.Fatal("expected allow for svc-x")
	}
	if !c.Allow(entry("svc-y")) {
		t.Fatal("expected allow for svc-y")
	}
	if c.Allow(entry("svc-x")) {
		t.Fatal("expected drop for svc-x after limit")
	}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}
