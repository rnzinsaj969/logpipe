// Package collect accumulates log entries up to a maximum count or until
// flushed, returning them as a snapshot slice.
package collect

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Collector buffers log entries in memory.
type Collector struct {
	mu      sync.Mutex
	entries []reader.LogEntry
	max     int
}

// New returns a Collector that holds at most max entries.
// When max is reached, the oldest entry is evicted (ring behaviour).
func New(max int) (*Collector, error) {
	if max <= 0 {
		return nil, errors.New("collect: max must be greater than zero")
	}
	return &Collector{max: max}, nil
}

// Add appends an entry to the collector.
// If the buffer is full the oldest entry is dropped.
func (c *Collector) Add(e reader.LogEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.max {
		c.entries = c.entries[1:]
	}
	c.entries = append(c.entries, e)
}

// Flush returns a copy of all buffered entries and resets the buffer.
func (c *Collector) Flush() []reader.LogEntry {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]reader.LogEntry, len(c.entries))
	copy(out, c.entries)
	c.entries = c.entries[:0]
	return out
}

// Len returns the current number of buffered entries.
func (c *Collector) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
