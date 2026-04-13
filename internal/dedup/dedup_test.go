package dedup

import (
	"sync"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestIsDuplicateFirstSeenReturnsFalse(t *testing.T) {
	d := New(5 * time.Second)
	e := Entry{Service: "svc", Message: "hello", Level: "info"}
	if d.IsDuplicate(e) {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicateWithinWindowReturnsTrue(t *testing.T) {
	now := time.Now()
	d := New(5 * time.Second)
	d.now = fixedClock(now)

	e := Entry{Service: "svc", Message: "hello", Level: "info"}
	d.IsDuplicate(e)

	if !d.IsDuplicate(e) {
		t.Fatal("second occurrence within window should be a duplicate")
	}
}

func TestIsDuplicateAfterWindowReturnsFalse(t *testing.T) {
	now := time.Now()
	d := New(2 * time.Second)
	d.now = fixedClock(now)

	e := Entry{Service: "svc", Message: "msg", Level: "warn"}
	d.IsDuplicate(e)

	// Advance clock beyond window
	d.now = fixedClock(now.Add(3 * time.Second))
	if d.IsDuplicate(e) {
		t.Fatal("occurrence after window should not be a duplicate")
	}
}

func TestEvictRemovesStaleEntries(t *testing.T) {
	now := time.Now()
	d := New(1 * time.Second)
	d.now = fixedClock(now)

	d.IsDuplicate(Entry{Service: "a", Message: "m1", Level: "info"})
	d.IsDuplicate(Entry{Service: "b", Message: "m2", Level: "error"})

	d.now = fixedClock(now.Add(2 * time.Second))
	d.Evict()

	if d.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", d.Len())
	}
}

func TestConcurrentIsDuplicate(t *testing.T) {
	d := New(10 * time.Second)
	e := Entry{Service: "svc", Message: "concurrent", Level: "debug"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.IsDuplicate(e)
		}()
	}
	wg.Wait()

	if d.Len() != 1 {
		t.Fatalf("expected 1 tracked entry, got %d", d.Len())
	}
}
