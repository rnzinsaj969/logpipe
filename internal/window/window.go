package window

import (
	"sync"
	"time"
)

// Entry holds a log entry timestamp used for windowed counting.
type Entry struct {
	At time.Time
}

// Counter tracks how many events occurred within a sliding time window.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	entries  []time.Time
	clock    func() time.Time
}

// New returns a Counter with the given sliding window duration.
func New(window time.Duration) *Counter {
	return NewWithClock(window, time.Now)
}

// NewWithClock returns a Counter using the provided clock function.
// Useful for deterministic testing.
func NewWithClock(window time.Duration, clock func() time.Time) *Counter {
	return &Counter{
		window: window,
		clock:  clock,
	}
}

// Add records a new event at the current clock time.
func (c *Counter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.entries = append(c.entries, now)
	c.evict(now)
}

// Count returns the number of events within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(c.clock())
	return len(c.entries)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = c.entries[:0]
}

// evict removes entries older than the window. Must be called with mu held.
func (c *Counter) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.entries) && c.entries[i].Before(cutoff) {
		i++
	}
	c.entries = c.entries[i:]
}
