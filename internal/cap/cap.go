// Package cap limits the number of log entries emitted per service within a
// sliding time window, dropping entries that exceed the configured ceiling.
package cap

import (
	"fmt"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options configures the Capper.
type Options struct {
	// Max is the maximum number of entries allowed per service per Window.
	Max int
	// Window is the duration of the sliding window.
	Window time.Duration
}

type bucket struct {
	count int
	start time.Time
}

// Capper drops log entries that exceed a per-service rate ceiling.
type Capper struct {
	opts    Options
	mu      sync.Mutex
	buckets map[string]*bucket
	now     func() time.Time
}

// New returns a Capper with the given options.
func New(opts Options) (*Capper, error) {
	if opts.Max <= 0 {
		return nil, fmt.Errorf("cap: Max must be greater than zero")
	}
	if opts.Window <= 0 {
		return nil, fmt.Errorf("cap: Window must be greater than zero")
	}
	return &Capper{
		opts:    opts,
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}, nil
}

// Allow returns true if the entry is within the cap for its service.
func (c *Capper) Allow(e reader.LogEntry) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	b, ok := c.buckets[e.Service]
	if !ok || now.Sub(b.start) >= c.opts.Window {
		c.buckets[e.Service] = &bucket{count: 1, start: now}
		return true
	}
	if b.count >= c.opts.Max {
		return false
	}
	b.count++
	return true
}
