package watermark

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Watermark tracks the highest timestamp seen across a stream of log entries.
// It is safe for concurrent use.
type Watermark struct {
	mu   sync.RWMutex
	high time.Time
}

// New returns a new Watermark with a zero high-water mark.
func New() *Watermark {
	return &Watermark{}
}

// Advance updates the high-water mark if the entry's timestamp is later than
// the current mark. It returns an error if the entry is nil.
func (w *Watermark) Advance(entry *reader.LogEntry) error {
	if entry == nil {
		return errors.New("watermark: nil entry")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if entry.Timestamp.After(w.high) {
		w.high = entry.Timestamp
	}
	return nil
}

// High returns the current high-water mark.
func (w *Watermark) High() time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.high
}

// Behind reports whether the given entry's timestamp is strictly before the
// current high-water mark, i.e. the entry is late / out-of-order.
func (w *Watermark) Behind(entry *reader.LogEntry) bool {
	if entry == nil {
		return false
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return entry.Timestamp.Before(w.high)
}

// Reset sets the high-water mark back to the zero time.
func (w *Watermark) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.high = time.Time{}
}
