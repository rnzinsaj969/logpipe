// Package counter provides a per-service log entry counter with
// snapshot and reset capabilities.
package counter

import "sync"

// Counter tracks the number of log entries seen per service.
type Counter struct {
	mu     sync.Mutex
	counts map[string]int64
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{counts: make(map[string]int64)}
}

// Inc increments the count for the given service by one.
func (c *Counter) Inc(service string) {
	c.mu.Lock()
	c.counts[service]++
	c.mu.Unlock()
}

// Add increments the count for the given service by n.
func (c *Counter) Add(service string, n int64) {
	if n <= 0 {
		return
	}
	c.mu.Lock()
	c.counts[service] += n
	c.mu.Unlock()
}

// Get returns the current count for the given service.
func (c *Counter) Get(service string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[service]
}

// Snapshot returns a copy of all current counts.
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Reset zeroes all counts.
func (c *Counter) Reset() {
	c.mu.Lock()
	c.counts = make(map[string]int64)
	c.mu.Unlock()
}
