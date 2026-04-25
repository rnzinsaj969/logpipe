// Package ceiling enforces a per-service maximum entry count within a
// rolling time window, dropping entries that exceed the configured cap.
package ceiling

import (
	"fmt"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Ceiling drops log entries for a service once the configured maximum
// count within the rolling window has been reached.
type Ceiling struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	counts  map[string][]time.Time
	nowFunc func() time.Time
}

// New creates a Ceiling that allows at most max entries per service
// within the given rolling window duration.
func New(max int, window time.Duration) (*Ceiling, error) {
	if max <= 0 {
		return nil, fmt.Errorf("ceiling: max must be positive, got %d", max)
	}
	if window <= 0 {
		return nil, fmt.Errorf("ceiling: window must be positive, got %s", window)
	}
	return &Ceiling{
		max:     max,
		window:  window,
		counts:  make(map[string][]time.Time),
		nowFunc: time.Now,
	}, nil
}

// Allow returns true when the entry is within the per-service limit and
// records the observation. It returns false when the cap is exceeded.
func (c *Ceiling) Allow(e reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.nowFunc()
	cutoff := now.Add(-c.window)

	times := c.counts[e.Service]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= c.max {
		c.counts[e.Service] = filtered
		return false
	}

	c.counts[e.Service] = append(filtered, now)
	return true
}

// Reset clears all per-service counters.
func (c *Ceiling) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string][]time.Time)
}
