package window_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/window"
)

func TestWindowRollingBehaviour(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	c := window.NewWithClock(3*time.Second, clock)

	for i := 0; i < 5; i++ {
		c.Add()
	}
	if got := c.Count(); got != 5 {
		t.Fatalf("step1: expected 5, got %d", got)
	}

	now = now.Add(4 * time.Second) // all entries expire
	if got := c.Count(); got != 0 {
		t.Fatalf("step2: expected 0, got %d", got)
	}

	for i := 0; i < 3; i++ {
		c.Add()
	}
	if got := c.Count(); got != 3 {
		t.Fatalf("step3: expected 3, got %d", got)
	}
}

func TestWindowResetThenRefill(t *testing.T) {
	c := window.New(10 * time.Second)
	c.Add()
	c.Add()
	c.Reset()
	c.Add()
	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1 after reset+add, got %d", got)
	}
}
