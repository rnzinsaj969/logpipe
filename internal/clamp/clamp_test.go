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

func (f *fixedClock) now() time.Time { return f.t }

func TestAllowWithinMax(t *testing.T) {
	clk := &fixedClock{t: time.Unix(0, 0)}
	l, err := newWithClock(3, time.Second, clk.now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !l.Allow(entry("svc")) {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllowExceedsMax(t *testing.T) {
	clk := &fixedClock{t: time.Unix(0, 0)}
	l, _ := newWithClock(2, time.Second, clk.now)
	l.Allow(entry("svc"))
	l.Allow(entry("svc"))
	if l.Allow(entry("svc")) {
		t.Fatal("expected Allow=false after exceeding max")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	clk := &fixedClock{t: time.Unix(0, 0)}
	l, _ := newWithClock(1, time.Second, clk.now)
	l.Allow(entry("svc"))
	if l.Allow(entry("svc")) {
		t.Fatal("expected Allow=false within window")
	}
	clk.t = clk.t.Add(2 * time.Second)
	if !l.Allow(entry("svc")) {
		t.Fatal("expected Allow=true after window reset")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	clk := &fixedClock{t: time.Unix(0, 0)}
	l, _ := newWithClock(1, time.Second, clk.now)
	if !l.Allow(entry("a")) {
		t.Fatal("expected true for service a")
	}
	if !l.Allow(entry("b")) {
		t.Fatal("expected true for service b")
	}
	if l.Allow(entry("a")) {
		t.Fatal("expected false for service a after cap")
	}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(10, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}
