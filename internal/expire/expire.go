// Package expire provides a processor that drops log entries older than a
// configurable maximum age.
package expire

import (
	"errors"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Processor drops entries whose timestamp is older than MaxAge.
type Processor struct {
	maxAge time.Duration
	clock  func() time.Time
}

// Options configures the Processor.
type Options struct {
	MaxAge time.Duration
	// Clock is optional; defaults to time.Now.
	Clock func() time.Time
}

// New returns a Processor that discards entries older than opts.MaxAge.
func New(opts Options) (*Processor, error) {
	if opts.MaxAge <= 0 {
		return nil, errors.New("expire: MaxAge must be positive")
	}
	clock := opts.Clock
	if clock == nil {
		clock = time.Now
	}
	return &Processor{maxAge: opts.MaxAge, clock: clock}, nil
}

// Apply returns false when the entry is expired, true otherwise.
func (p *Processor) Apply(e reader.LogEntry) bool {
	if e.Timestamp.IsZero() {
		return true
	}
	return p.clock().Sub(e.Timestamp) <= p.maxAge
}

// Filter returns only non-expired entries from the provided slice.
func (p *Processor) Filter(entries []reader.LogEntry) []reader.LogEntry {
	out := entries[:0:0]
	for _, e := range entries {
		if p.Apply(e) {
			out = append(out, e)
		}
	}
	return out
}
