// Package ceiling provides a per-service log entry ceiling that hard-caps
// the total number of entries accepted over the lifetime of the processor.
// Once the ceiling is reached for a given service, all further entries from
// that service are dropped until Reset is called.
package ceiling

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Ceiling drops entries for a service once its lifetime count exceeds Max.
type Ceiling struct {
	mu     sync.Mutex
	max    int
	counts map[string]int
}

// New creates a Ceiling with the given maximum per-service entry count.
// Max must be greater than zero.
func New(max int) (*Ceiling, error) {
	if max <= 0 {
		return nil, errors.New("ceiling: max must be greater than zero")
	}
	return &Ceiling{
		max:    max,
		counts: make(map[string]int),
	}, nil
}

// Allow reports whether the entry should be passed downstream.
// It increments the counter for entry.Service and returns false once the
// ceiling has been reached.
func (c *Ceiling) Allow(entry reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := entry.Service
	if key == "" {
		key = "__default__"
	}

	if c.counts[key] >= c.max {
		return false
	}
	c.counts[key]++
	return true
}

// Reset clears all counters, allowing entries to flow again.
func (c *Ceiling) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]int)
}

// Snapshot returns a copy of the current per-service counts.
func (c *Ceiling) Snapshot() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}
