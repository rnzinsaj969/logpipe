// Package ceiling enforces a per-service maximum entry count within a
// rolling time window, dropping entries that exceed the configured limit.
package ceiling

import (
	"fmt"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Ceiling drops log entries for a service once the configured maximum
// count has been reached within the current window.
type Ceiling struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	counts  map[string]int
	resets  map[string]time.Time
	nowFunc func() time.Time
}

// New creates a Ceiling that allows at most max entries per service
// within window. It returns an error when max is not positive.
func New(max int, window time.Duration) (*Ceiling, error) {
	if max <= 0 {
		return nil, fmt.Errorf("ceiling: max must be positive, got %d", max)
	}
	return &Ceiling{
		max:     max,
		window:  window,
		counts:  make(map[string]int),
		resets:  make(map[string]time.Time),
		nowFunc: time.Now,
	}, nil
}

// Allow returns true when the entry is within the per-service limit for
// the current window, and false when the ceiling has been reached.
func (c *Ceiling) Allow(e reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.nowFunc()
	svc := e.Service

	if reset, ok := c.resets[svc]; !ok || now.After(reset) {
		c.counts[svc] = 0
		c.resets[svc] = now.Add(c.window)
	}

	if c.counts[svc] >= c.max {
		return false
	}
	c.counts[svc]++
	return true
}
