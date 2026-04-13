package dedup

import (
	"sync"
	"time"
)

// Entry represents a minimal log entry used for deduplication keying.
type Entry struct {
	Service string
	Message string
	Level   string
}

// Deduplicator suppresses repeated log entries within a configurable window.
type Deduplicator struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[Entry]time.Time
	now    func() time.Time
}

// New creates a Deduplicator that suppresses duplicate entries seen within window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window: window,
		seen:   make(map[Entry]time.Time),
		now:    time.Now,
	}
}

// IsDuplicate returns true if an identical entry was seen within the dedup window.
// If it is not a duplicate, the entry is recorded and false is returned.
func (d *Deduplicator) IsDuplicate(e Entry) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if last, ok := d.seen[e]; ok && now.Sub(last) < d.window {
		return true
	}

	d.seen[e] = now
	return false
}

// Evict removes stale entries that are outside the dedup window.
// Call periodically to prevent unbounded memory growth.
func (d *Deduplicator) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, last := range d.seen {
		if now.Sub(last) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of entries currently tracked.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
