// Package dedupe provides field-level deduplication for log entries,
// removing repeated key-value pairs across a sliding window of entries.
package dedupe

import (
	"sync"
	"time"

	"logpipe/internal/reader"
)

// Options configures the Deduplicator.
type Options struct {
	// Fields is the list of LogEntry fields to compare for deduplication.
	// Supported values: "message", "level", "service".
	Fields []string
	// Window is the duration within which identical entries are suppressed.
	Window time.Duration
}

// DefaultOptions returns sensible defaults for deduplication.
func DefaultOptions() Options {
	return Options{
		Fields: []string{"message", "service"},
		Window: 5 * time.Second,
	}
}

// Deduplicator suppresses log entries whose key fields match a recently seen entry.
type Deduplicator struct {
	opts  Options
	mu    sync.Mutex
	seen  map[string]time.Time
	clock func() time.Time
}

// New creates a Deduplicator with the given options.
func New(opts Options) *Deduplicator {
	return &Deduplicator{
		opts:  opts,
		seen:  make(map[string]time.Time),
		clock: time.Now,
	}
}

// IsDuplicate returns true if an equivalent entry was seen within the window.
func (d *Deduplicator) IsDuplicate(e reader.LogEntry) bool {
	key := d.buildKey(e)
	now := d.clock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.seen[key]; ok && now.Sub(t) < d.opts.Window {
		return true
	}
	d.seen[key] = now
	return false
}

// Evict removes entries older than the configured window.
func (d *Deduplicator) Evict() {
	now := d.clock()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, t := range d.seen {
		if now.Sub(t) >= d.opts.Window {
			delete(d.seen, k)
		}
	}
}

// buildKey constructs a string key from the configured fields of the entry.
func (d *Deduplicator) buildKey(e reader.LogEntry) string {
	var key string
	for _, f := range d.opts.Fields {
		switch f {
		case "message":
			key += "msg=" + e.Message + ";"
		case "level":
			key += "lvl=" + e.Level + ";"
		case "service":
			key += "svc=" + e.Service + ";"
		}
	}
	return key
}
