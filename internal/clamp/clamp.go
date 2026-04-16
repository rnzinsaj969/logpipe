// Package clamp provides a processor that enforces minimum and maximum
// entry counts per service within a sliding time window.
package clamp

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options configures the Clamper.
type Options struct {
	Min      int
	Max      int
	Window   time.Duration
	NowFunc  func() time.Time
}

// Clamper tracks per-service entry counts and drops entries that exceed Max
// or have not yet reached Min within the current window.
type Clamper struct {
	opts    Options
	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New returns a Clamper. Max must be >= Min >= 0 and Window > 0.
func New(opts Options) (*Clamper, error) {
	if opts.Window <= 0 {
		return nil, errors.New("clamp: window must be positive")
	}
	if opts.Max < opts.Min {
		return nil, errors.New("clamp: max must be >= min")
	}
	if opts.NowFunc == nil {
		opts.NowFunc = time.Now
	}
	return &Clamper{opts: opts, buckets: make(map[string]*bucket)}, nil
}

// Allow returns true when the entry should be forwarded.
func (c *Clamper) Allow(e reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.opts.NowFunc()
	b, ok := c.buckets[e.Service]
	if !ok || now.After(b.windowEnd) {
		b = &bucket{windowEnd: now.Add(c.opts.Window)}
		c.buckets[e.Service] = b
	}
	b.count++
	if c.opts.Max > 0 && b.count > c.opts.Max {
		return false
	}
	return true
}
