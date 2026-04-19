// Package debounce suppresses repeated log entries within a configurable
// quiet period, emitting only the first occurrence until the window expires.
package debounce

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// DefaultWindow is the default debounce duration.
const DefaultWindow = 5 * time.Second

type clock func() time.Time

type entry struct {
	first time.Time
}

// Debouncer suppresses entries whose (service, message) key has been seen
// within the active window.
type Debouncer struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]entry
	now    clock
}

// New returns a Debouncer with the given quiet window.
// Returns an error if window is not positive.
func New(window time.Duration) (*Debouncer, error) {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, c clock) (*Debouncer, error) {
	if window <= 0 {
		return nil, errors.New("debounce: window must be positive")
	}
	return &Debouncer{
		window: window,
		seen:   make(map[string]entry),
		now:    c,
	}, nil
}

// Allow returns true if the entry should be forwarded (first occurrence or
// window expired), false if it should be suppressed.
func (d *Debouncer) Allow(e reader.LogEntry) bool {
	key := e.Service + "\x00" + e.Message
	now := d.now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if rec, ok := d.seen[key]; ok {
		if now.Sub(rec.first) < d.window {
			return false
		}
	}
	d.seen[key] = entry{first: now}
	return true
}

// Evict removes stale entries whose window has expired.
func (d *Debouncer) Evict() {
	now := d.now()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, rec := range d.seen {
		if now.Sub(rec.first) >= d.window {
			delete(d.seen, k)
		}
	}
}
