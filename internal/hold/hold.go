// Package hold provides a conditional entry buffer that accumulates log
// entries until a release predicate is satisfied, then flushes them all at
// once. Entries can be discarded instead of released by calling Discard.
package hold

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Predicate returns true when the held entries should be released.
type Predicate func(e reader.LogEntry) bool

// Holder buffers entries and releases them when a predicate fires.
type Holder struct {
	mu      sync.Mutex
	buf     []reader.LogEntry
	max     int
	pred    Predicate
}

// New creates a Holder with the given capacity and release predicate.
// max is the maximum number of entries to retain; older entries are dropped
// when the buffer is full. pred must not be nil.
func New(max int, pred Predicate) (*Holder, error) {
	if max <= 0 {
		return nil, errors.New("hold: max must be greater than zero")
	}
	if pred == nil {
		return nil, errors.New("hold: predicate must not be nil")
	}
	return &Holder{buf: make([]reader.LogEntry, 0, max), max: max, pred: pred}, nil
}

// Add buffers e. If the release predicate returns true for e, all buffered
// entries (including e) are returned and the buffer is reset.
// If the buffer is full the oldest entry is evicted to make room.
func (h *Holder) Add(e reader.LogEntry) ([]reader.LogEntry, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.buf) >= h.max {
		h.buf = h.buf[1:]
	}
	h.buf = append(h.buf, e)

	if h.pred(e) {
		out := make([]reader.LogEntry, len(h.buf))
		copy(out, h.buf)
		h.buf = h.buf[:0]
		return out, true
	}
	return nil, false
}

// Discard drops all buffered entries without releasing them.
func (h *Holder) Discard() {
	h.mu.Lock()
	h.buf = h.buf[:0]
	h.mu.Unlock()
}

// Len returns the number of currently buffered entries.
func (h *Holder) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.buf)
}
