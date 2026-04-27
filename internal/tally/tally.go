// Package tally provides a per-field frequency counter for log entries.
// It tracks how often each unique value appears for a given field across
// a stream of entries and exposes a snapshot of the current counts.
package tally

import (
	"errors"
	"sync"

	"github.com/yourorg/logpipe/internal/reader"
)

// Tally counts occurrences of each unique value for a named field.
type Tally struct {
	mu    sync.Mutex
	field string
	counts map[string]int
}

// New creates a Tally that counts values for the given field name.
// field may be "level", "service", or any key present in Extra.
// Returns an error if field is empty.
func New(field string) (*Tally, error) {
	if field == "" {
		return nil, errors.New("tally: field must not be empty")
	}
	return &Tally{
		field:  field,
		counts: make(map[string]int),
	}, nil
}

// Add records the value of the configured field from entry.
// Zero-value strings are counted under the empty-string key.
func (t *Tally) Add(entry reader.LogEntry) {
	val := t.valueOf(entry)
	t.mu.Lock()
	t.counts[val]++
	t.mu.Unlock()
}

// Snapshot returns a copy of the current counts map.
func (t *Tally) Snapshot() map[string]int {
	t.mu.Lock()
	out := make(map[string]int, len(t.counts))
	for k, v := range t.counts {
		out[k] = v
	}
	t.mu.Unlock()
	return out
}

// Reset clears all accumulated counts.
func (t *Tally) Reset() {
	t.mu.Lock()
	t.counts = make(map[string]int)
	t.mu.Unlock()
}

func (t *Tally) valueOf(entry reader.LogEntry) string {
	switch t.field {
	case "level":
		return entry.Level
	case "service":
		return entry.Service
	case "message":
		return entry.Message
	default:
		if entry.Extra != nil {
			if v, ok := entry.Extra[t.field]; ok {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
		return ""
	}
}
