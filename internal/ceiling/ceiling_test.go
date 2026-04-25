package ceiling

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "test", Level: "info"}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAllowWithinLimit(t *testing.T) {
	c, err := New(3, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	c, err := New(2, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny on third call")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	c, _ := New(1, time.Minute)
	if !c.Allow(entry("a")) {
		t.Fatal("expected allow for service a")
	}
	if !c.Allow(entry("b")) {
		t.Fatal("expected allow for service b")
	}
	if c.Allow(entry("a")) {
		t.Fatal("expected deny for service a second call")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	c, _ := New(1, 50*time.Millisecond)
	c.nowFunc = func() time.Time { return time.Now().Add(-100 * time.Millisecond) }
	c.Allow(entry("svc"))
	c.nowFunc = time.Now
	if !c.Allow(entry("svc")) {
		t.Fatal("expected allow after window expiry")
	}
}

func TestResetClearsState(t *testing.T) {
	c, _ := New(1, time.Minute)
	c.Allow(entry("svc"))
	c.Reset()
	if !c.Allow(entry("svc")) {
		t.Fatal("expected allow after reset")
	}
}
