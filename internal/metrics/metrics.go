package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter tracks a monotonically increasing integer value.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

// Value returns the current counter value.
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Snapshot holds a point-in-time view of all tracked metrics.
type Snapshot struct {
	EntriesRead    uint64
	EntriesMatched uint64
	EntriesDropped uint64
	ParseErrors    uint64
}

// Metrics aggregates runtime counters for the logpipe pipeline.
type Metrics struct {
	mu             sync.RWMutex
	EntriesRead    Counter
	EntriesMatched Counter
	EntriesDropped Counter
	ParseErrors    Counter
}

// New returns a new, zeroed Metrics instance.
func New() *Metrics {
	return &Metrics{}
}

// Snapshot returns a consistent point-in-time copy of all counters.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Snapshot{
		EntriesRead:    m.EntriesRead.Value(),
		EntriesMatched: m.EntriesMatched.Value(),
		EntriesDropped: m.EntriesDropped.Value(),
		ParseErrors:    m.ParseErrors.Value(),
	}
}
