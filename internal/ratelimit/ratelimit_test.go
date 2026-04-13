package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAllowConsumesTokens(t *testing.T) {
	l := New(3)
	// Burst capacity is 3; first three calls should succeed.
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	// Fourth call should be denied.
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
}

func TestAllowRefillsOverTime(t *testing.T) {
	l := New(10)
	// Drain all tokens.
	for i := 0; i < 10; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected bucket to be empty")
	}

	// Advance the internal clock by 500ms — should add ~5 tokens.
	base := time.Now()
	l.clock = func() time.Time { return base.Add(500 * time.Millisecond) }
	l.lastTick = base

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed < 4 || allowed > 6 {
		t.Fatalf("expected ~5 tokens after 500ms refill, got %d", allowed)
	}
}

func TestReset(t *testing.T) {
	l := New(5)
	for i := 0; i < 5; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected bucket empty before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected Allow()=true after reset")
	}
}

func TestConcurrentAllow(t *testing.T) {
	const limit = 100
	l := New(limit)

	var wg sync.WaitGroup
	var allowed atomic.Int64

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got > limit {
		t.Fatalf("concurrent Allow() granted %d > burst limit %d", got, limit)
	}
}
