// Package clamp limits the rate of log entries per service within a sliding
// time window, dropping entries that exceed the configured maximum.
package clamp

import (
	"fmt"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// clock is a small interface so tests can inject a fake time source.
type clock func() time.Time

// entry tracks the count of log entries seen within the current window.
type entry struct {
	count    int
	windowAt time.Time
}

// Clamp drops log entries for a given service once they exceed Max within
// the configured Window duration.
type Clamp struct {
	mu     sync.Mutex
	max    int
	win    time.Duration
	clock  clock
	bucket map[string]*entry
}

// New returns a Clamp that allows at most max entries per service per window.
// max must be positive and window must be greater than zero.
func New(max int, window time.Duration) (*Clamp, error) {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, c clock) (*Clamp, error) {
	if max <= 0 {
		return nil, fmt.Errorf("clamp: max must be positive, got %d", max)
	}
	if window <= 0 {
		return nil, fmt.Errorf("clamp: window must be positive, got %s", window)
	}
	return &Clamp{
		max:    max,
		win:    window,
		clock:  c,
		bucket: make(map[string]*entry),
	}, nil
}

// Allow returns true when the entry is within the rate limit for its service
// and false when it should be dropped.
func (c *Clamp) Allow(e reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.clock()
	svc := e.Service

	st, ok := c.bucket[svc]
	if !ok || now.Sub(st.windowAt) >= c.win {
		c.bucket[svc] = &entry{count: 1, windowAt: now}
		return true
	}
	st.count++
	return st.count <= c.max
}
