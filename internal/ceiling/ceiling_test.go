package ceiling

import (
	"testing"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(service string) reader.LogEntry {
	return reader.LogEntry{Service: service, Message: "test", Level: "info"}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAllowWithinLimit(t *testing.T) {
	c, _ := New(3)
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected Allow=true on iteration %d", i)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	c, _ := New(2)
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected Allow=false after ceiling reached")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	c, _ := New(1)
	if !c.Allow(entry("a")) {
		t.Fatal("expected Allow=true for service a")
	}
	if !c.Allow(entry("b")) {
		t.Fatal("expected Allow=true for service b")
	}
	if c.Allow(entry("a")) {
		t.Fatal("expected Allow=false for service a after ceiling")
	}
}

func TestResetRestoresAllowance(t *testing.T) {
	c, _ := New(1)
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected Allow=false before reset")
	}
	c.Reset()
	if !c.Allow(entry("svc")) {
		t.Fatal("expected Allow=true after reset")
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	c, _ := New(5)
	c.Allow(entry("x"))
	c.Allow(entry("x"))
	snap := c.Snapshot()
	snap["x"] = 999
	if c.Snapshot()["x"] != 2 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestEmptyServiceUsesDefault(t *testing.T) {
	c, _ := New(1)
	if !c.Allow(entry("")) {
		t.Fatal("expected Allow=true for empty service")
	}
	if c.Allow(entry("")) {
		t.Fatal("expected Allow=false for empty service after ceiling")
	}
}
