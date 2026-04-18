// Package group provides log entry grouping by a specified field.
package group

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Grouper accumulates log entries bucketed by a field value.
type Grouper struct {
	mu    sync.Mutex
	field string
	buckets map[string][]reader.LogEntry
}

// New returns a Grouper that buckets entries by field.
// field may be "level", "service", or any key in Extra.
func New(field string) (*Grouper, error) {
	if field == "" {
		return nil, errors.New("group: field must not be empty")
	}
	return &Grouper{
		field:   field,
		buckets: make(map[string][]reader.LogEntry),
	}, nil
}

// Add places entry into the appropriate bucket.
func (g *Grouper) Add(entry reader.LogEntry) {
	key := g.keyFor(entry)
	g.mu.Lock()
	g.buckets[key] = append(g.buckets[key], entry)
	g.mu.Unlock()
}

// Snapshot returns a copy of the current buckets and resets state.
func (g *Grouper) Snapshot() map[string][]reader.LogEntry {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make(map[string][]reader.LogEntry, len(g.buckets))
	for k, v := range g.buckets {
		cp := make([]reader.LogEntry, len(v))
		copy(cp, v)
		out[k] = cp
	}
	g.buckets = make(map[string][]reader.LogEntry)
	return out
}

func (g *Grouper) keyFor(e reader.LogEntry) string {
	switch g.field {
	case "level":
		return e.Level
	case "service":
		return e.Service
	default:
		if e.Extra != nil {
			if v, ok := e.Extra[g.field]; ok {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
		return ""
	}
}
