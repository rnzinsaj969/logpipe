// Package stamp provides a processor that overwrites or sets timestamp
// fields on log entries using a configurable clock source.
package stamp

import (
	"errors"
	"time"

	"logpipe/internal/reader"
)

// Options controls Stamper behaviour.
type Options struct {
	// OverwriteExisting replaces a non-zero timestamp when true.
	OverwriteExisting bool
	// Clock is used to obtain the current time. Defaults to time.Now.
	Clock func() time.Time
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		OverwriteExisting: false,
		Clock:             time.Now,
	}
}

// Stamper applies a timestamp to log entries.
type Stamper struct {
	opts Options
}

// New creates a Stamper. Returns an error if Clock is nil.
func New(opts Options) (*Stamper, error) {
	if opts.Clock == nil {
		return nil, errors.New("stamp: Clock must not be nil")
	}
	return &Stamper{opts: opts}, nil
}

// Apply sets the Timestamp on e according to the configured options.
// It returns a new entry and never mutates the original.
func (s *Stamper) Apply(e reader.LogEntry) reader.LogEntry {
	if !e.Timestamp.IsZero() && !s.opts.OverwriteExisting {
		return e
	}
	out := e
	out.Timestamp = s.opts.Clock()
	return out
}
