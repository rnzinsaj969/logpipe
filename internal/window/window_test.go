package window

import (
	"sync"
	"testing"
	"time"
)

func TestCountIncreasesOnAdd(t *testing.T) {
	c := New(5 * time.Second)
	c.Add()
	c.Add()
	if got := c.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCountExcludesExpiredEntries(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	c := NewWithClock(2*time.Second, clock)

	c.Add() // recorded at 'now'

	now = now.Add(3 * time.Second) // advance past window
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCountKeepsEntriesWithinWindow(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	c := NewWithClock(5*time.Second, clock)

	c.Add()
	now = now.Add(2 * time.Second)
	c.Add()

	if got := c.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestResetClearsAll(t *testing.T) {
	c := New(10 * time.Second)
	c.Add()
	c.Add()
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestConcurrentAdd(t *testing.T) {
	c := New(10 * time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add()
		}()
	}
	wg.Wait()
	if got := c.Count(); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}

func TestPartialEviction(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	c := NewWithClock(4*time.Second, clock)

	c.Add() // t=0
	now = now.Add(2 * time.Second)
	c.Add() // t=2
	now = now.Add(3 * time.Second) // window cutoff at t=1; first entry expires

	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
