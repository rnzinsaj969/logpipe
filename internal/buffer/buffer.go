package buffer

import (
	"sync"

	"github.com/yourorg/logpipe/internal/filter"
)

// RingBuffer holds a fixed-size circular buffer of log entries.
// Once full, the oldest entries are overwritten by new ones.
type RingBuffer struct {
	mu      sync.Mutex
	items   []filter.LogEntry
	cap     int
	head    int
	count   int
}

// New creates a new RingBuffer with the given capacity.
// Panics if capacity is less than 1.
func New(capacity int) *RingBuffer {
	if capacity < 1 {
		panic("buffer: capacity must be at least 1")
	}
	return &RingBuffer{
		items: make([]filter.LogEntry, capacity),
		cap:   capacity,
	}
}

// Push adds a log entry to the buffer. If the buffer is full,
// the oldest entry is silently overwritten.
func (r *RingBuffer) Push(entry filter.LogEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := (r.head + r.count) % r.cap
	r.items[index] = entry

	if r.count < r.cap {
		r.count++
	} else {
		// overwrite: advance head
		r.head = (r.head + 1) % r.cap
	}
}

// Drain returns all buffered entries in insertion order and resets the buffer.
func (r *RingBuffer) Drain() []filter.LogEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]filter.LogEntry, r.count)
	for i := 0; i < r.count; i++ {
		out[i] = r.items[(r.head+i)%r.cap]
	}
	r.head = 0
	r.count = 0
	return out
}

// Len returns the current number of entries in the buffer.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}
