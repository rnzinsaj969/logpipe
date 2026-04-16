package clamp_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/clamp"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "msg", Level: "info"}
}

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestAllowWithinMax(t *testing.T) {
	c, err := clamp.New(clamp.Options{Min: 0, Max: 3, Window: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if !c.Allow(entry("svc")) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsMax(t *testing.T) {
	c, _ := clamp.New(clamp.Options{Max: 2, Window: time.Second})
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny on third call")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	now := time.Now()
	calls := 0
	clock := func() time.Time {
		calls++
		if calls <= 4 {
			return now
		}
		return now.Add(2 * time.Second)
	}
	c, _ := clamp.New(clamp.Options{Max: 2, Window: time.Second, NowFunc: clock})
	c.Allow(entry("svc"))
	c.Allow(entry("svc"))
	if c.Allow(entry("svc")) {
		t.Fatal("expected deny before window reset")
	}
	if !c.Allow(entry("svc")) {
		t.Fatal("expected allow after window reset")
	}
}

func TestIndependentServices(t *testing.T) {
	c, _ := clamp.New(clamp.Options{Max: 1, Window: time.Second})
	if !c.Allow(entry("a")) {
		t.Fatal("expected allow for a")
	}
	if !c.Allow(entry("b")) {
		t.Fatal("expected allow for b")
	}
	if c.Allow(entry("a")) {
		t.Fatal("expected deny for a on second call")
	}
}

func TestInvalidOptions(t *testing.T) {
	_, err := clamp.New(clamp.Options{Min: 5, Max: 2, Window: time.Second})
	if err == nil {
		t.Fatal("expected error for max < min")
	}
	_, err = clamp.New(clamp.Options{Max: 1, Window: 0})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}
