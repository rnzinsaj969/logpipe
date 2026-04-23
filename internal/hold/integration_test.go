package hold_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/hold"
	"github.com/logpipe/logpipe/internal/reader"
)

func TestConcurrentAddIsSafe(t *testing.T) {
	h, _ := hold.New(100, func(e reader.LogEntry) bool { return e.Level == "error" })

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Add(entry("msg", "info"))
		}()
	}
	wg.Wait()
}

func TestHoldReleaseAndRefill(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return e.Level == "error" })

	h.Add(entry("a", "info"))
	out, released := h.Add(entry("b", "error"))
	if !released || len(out) != 2 {
		t.Fatalf("first release: want 2 entries released=true, got %d released=%v", len(out), released)
	}

	// Buffer should be empty; refill and release again.
	h.Add(entry("c", "warn"))
	out2, released2 := h.Add(entry("d", "error"))
	if !released2 || len(out2) != 2 {
		t.Fatalf("second release: want 2 entries released=true, got %d released=%v", len(out2), released2)
	}
}

func TestDiscardThenRelease(t *testing.T) {
	h, _ := hold.New(10, func(e reader.LogEntry) bool { return e.Level == "error" })

	h.Add(entry("a", "info"))
	h.Add(entry("b", "info"))
	h.Discard()

	out, released := h.Add(entry("c", "error"))
	if !released {
		t.Fatal("expected release after discard")
	}
	if len(out) != 1 {
		t.Fatalf("want 1 entry after discard+add, got %d", len(out))
	}
}
