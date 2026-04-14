package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/throttle"
)

func fixedClock(t time.Time) throttle.Clock {
	return func() time.Time { return t }
}

func TestAllowWithinLimit(t *testing.T) {
	now := time.Now()
	th := throttle.New(time.Second, 3, fixedClock(now))

	for i := 0; i < 3; i++ {
		if !th.Allow("svc") {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	now := time.Now()
	th := throttle.New(time.Second, 3, fixedClock(now))

	for i := 0; i < 3; i++ {
		th.Allow("svc")
	}
	if th.Allow("svc") {
		t.Fatal("expected Allow=false after limit exceeded")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	base := time.Now()
	th := throttle.New(time.Second, 2, fixedClock(base))

	th.Allow("svc")
	th.Allow("svc")

	// advance past the window
	future := base.Add(2 * time.Second)
	th2 := throttle.New(time.Second, 2, fixedClock(future))
	// use a fresh throttler at future time to simulate window slide
	_ = th2

	// rebuild with a mutable clock
	current := base
	clock := func() time.Time { return current }
	th3 := throttle.New(time.Second, 2, clock)
	th3.Allow("svc")
	th3.Allow("svc")
	current = base.Add(2 * time.Second)
	if !th3.Allow("svc") {
		t.Fatal("expected Allow=true after window has passed")
	}
}

func TestAllowIndependentServices(t *testing.T) {
	now := time.Now()
	th := throttle.New(time.Second, 1, fixedClock(now))

	if !th.Allow("svcA") {
		t.Fatal("svcA first call should be allowed")
	}
	if !th.Allow("svcB") {
		t.Fatal("svcB first call should be allowed independently")
	}
	if th.Allow("svcA") {
		t.Fatal("svcA second call should be throttled")
	}
}

func TestReset(t *testing.T) {
	now := time.Now()
	th := throttle.New(time.Second, 1, fixedClock(now))
	th.Allow("svc")
	th.Reset()
	if !th.Allow("svc") {
		t.Fatal("expected Allow=true after Reset")
	}
}

func TestConcurrentAllow(t *testing.T) {
	th := throttle.New(time.Second, 100, nil)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			th.Allow("svc")
		}()
	}
	wg.Wait()
}
