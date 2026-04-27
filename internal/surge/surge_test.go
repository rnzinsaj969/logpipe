package surge

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func entry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "msg", Level: "info"}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(0, 2.0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewInvalidMultipleReturnsError(t *testing.T) {
	_, err := New(time.Second, 1.0)
	if err == nil {
		t.Fatal("expected error for multiple <= 1.0")
	}
}

func TestSingleEntryNoSurge(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Minute, 2.0, fixedClock(now))
	if d.Record(entry("svc")) {
		t.Fatal("single entry should not trigger surge")
	}
}

func TestNoSurgeUnderThreshold(t *testing.T) {
	base := time.Now()
	calls := 0
	clock := func() time.Time {
		calls++
		// spread 10 events evenly over the window
		return base.Add(time.Duration(calls-1) * 6 * time.Second)
	}
	d, _ := newWithClock(time.Minute, 3.0, clock)
	e := entry("svc")
	var last bool
	for i := 0; i < 10; i++ {
		last = d.Record(e)
	}
	if last {
		t.Fatal("evenly distributed events should not trigger surge")
	}
}

func TestSurgeDetectedOnBurst(t *testing.T) {
	base := time.Now()
	call := 0
	clock := func() time.Time {
		call++
		if call <= 5 {
			// 5 events spread over first 50 s (baseline ~0.1/s)
			return base.Add(time.Duration(call-1) * 10 * time.Second)
		}
		// next 5 events crammed into the last 5 s of the window
		return base.Add(55*time.Second + time.Duration(call-6)*time.Second)
	}
	d, _ := newWithClock(time.Minute, 2.0, clock)
	e := entry("svc")
	var last bool
	for i := 0; i < 10; i++ {
		last = d.Record(e)
	}
	if !last {
		t.Fatal("burst in recent quarter should trigger surge")
	}
}

func TestResetClearsState(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Minute, 2.0, fixedClock(now))
	d.Record(entry("svc"))
	d.Reset()
	if d.Record(entry("svc")) {
		t.Fatal("after reset a single entry must not trigger surge")
	}
}

func TestIndependentServices(t *testing.T) {
	now := time.Now()
	d, _ := newWithClock(time.Minute, 2.0, fixedClock(now))
	d.Record(entry("alpha"))
	if d.Record(entry("beta")) {
		t.Fatal("beta should be independent of alpha")
	}
}
