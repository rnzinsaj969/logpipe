// Package merge provides utilities for combining fields from multiple log
// entries into a single consolidated entry.
package merge

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls how entries are merged.
type Options struct {
	// PreferFirst keeps the value from the first entry when both entries
	// define the same top-level field. When false the last entry wins.
	PreferFirst bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{PreferFirst: true}
}

// Merger merges a slice of log entries into one.
type Merger struct {
	opts Options
}

// New creates a Merger with the provided options.
func New(opts Options) (*Merger, error) {
	return &Merger{opts: opts}, nil
}

// Apply merges entries into a single LogEntry. The first entry in the slice
// provides the base values for Service, Level, Message and Timestamp.
// Extra fields from all entries are combined according to PreferFirst.
// At least one entry must be provided.
func (m *Merger) Apply(entries []reader.LogEntry) (reader.LogEntry, error) {
	if len(entries) == 0 {
		return reader.LogEntry{}, errors.New("merge: no entries provided")
	}

	base := entries[0]

	// Build a merged extra map.
	merged := make(map[string]any)

	sources := entries
	if m.opts.PreferFirst {
		// Iterate in reverse so that earlier entries overwrite later ones.
		sources = make([]reader.LogEntry, len(entries))
		for i, e := range entries {
			sources[len(entries)-1-i] = e
		}
	}

	for _, e := range sources {
		for k, v := range e.Extra {
			merged[k] = v
		}
	}

	if len(merged) == 0 {
		merged = nil
	}

	return reader.LogEntry{
		Service:   base.Service,
		Level:     base.Level,
		Message:   base.Message,
		Timestamp: base.Timestamp,
		Extra:     merged,
	}, nil
}
