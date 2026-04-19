// Package pivot groups log entries by a field and counts occurrences.
package pivot

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Table holds aggregated counts keyed by a field value.
type Table struct {
	mu     sync.Mutex
	field  string
	counts map[string]int
}

// New creates a Pivot Table that groups entries by the given field name.
// field may be "level", "service", or any key inside Extra.
func New(field string) (*Table, error) {
	if field == "" {
		return nil, errors.New("pivot: field must not be empty")
	}
	return &Table{field: field, counts: make(map[string]int)}, nil
}

// Add records one entry into the table.
func (t *Table) Add(e reader.LogEntry) {
	key := t.resolve(e)
	t.mu.Lock()
	t.counts[key]++
	t.mu.Unlock()
}

// Snapshot returns a copy of the current counts and resets the table.
func (t *Table) Snapshot() map[string]int {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]int, len(t.counts))
	for k, v := range t.counts {
		out[k] = v
	}
	t.counts = make(map[string]int)
	return out
}

// Len returns the number of distinct keys currently tracked.
func (t *Table) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.counts)
}

func (t *Table) resolve(e reader.LogEntry) string {
	switch t.field {
	case "level":
		return e.Level
	case "service":
		return e.Service
	default:
		if e.Extra != nil {
			if v, ok := e.Extra[t.field]; ok {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
		return ""
	}
}
