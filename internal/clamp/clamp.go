// Package clamp limits the rate of log entries per service within a sliding
// time window, dropping entries that exceed the configured maximum.
package clamp

import (
	"fmt"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// clock is a narrow interface for obtaining the current time, allowing
// deterministic testing without real-time dependencies.
type clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Clamp drops log entries for a given service once the count within the
// sliding window exceeds Max.
type Clamp struct {
	max    int
	win    time.Duration
	clock  clock
	mu     sync.Mutex
	bucket map[string][]time.Time
}

// New returns a Clamp that allows at most max entries per service within win.
// max must be positive and win must be greater than zero.
func New(max int, win time.Duration) (*Clamp, error) {
	return newWithClock(max, win, realClock{})
}

func newWithClock(max int, win time.Duration, c clock) (*Clamp, error) {
	if max <= 0 {
		return nil, fmt.Errorf("clamp: max must be positive, got %d", max)
	}
	if win <= 0 {
		return nil, fmt.Errorf("clamp: window must be positive, got %s", win)
	}
	return &Clamp{
		max:    max,
		win:    win,
		clock:  c,
		bucket: make(map[string][]time.Time),
	}, nil
}

// Allow reports whether the entry should be forwarded. Entries that push the
// per-service count above Max within the current window are rejected.
func (c *Clamp) Allow(e reader.LogEntry) bool {
	now := c.clock.Now()
	cutoff := now.Add(-c.win)

	c.mu.Lock()
	defer c.mu.Unlock()

	times := c.bucket[e.Service]
	var kept []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= c.max {
		c.bucket[e.Service] = kept
		return false
	}

	c.bucket[e.Service] = append(kept, now)
	return true
}
