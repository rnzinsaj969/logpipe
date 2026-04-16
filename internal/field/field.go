// Package field provides utilities for copying, renaming, and deleting
// fields within a log entry's Extra map.
package field

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Op describes a single field operation.
type Op struct {
	// Action is one of "copy", "rename", "delete".
	Action string
	// From is the source key (unused for delete).
	From string
	// To is the destination key (used for copy and rename).
	To string
}

// Processor applies a sequence of field Ops to log entries.
type Processor struct {
	ops []Op
}

// New returns a Processor for the given ops. It returns an error if any op
// is invalid.
func New(ops []Op) (*Processor, error) {
	for _, o := range ops {
		switch o.Action {
		case "copy", "rename":
			if o.From == "" || o.To == "" {
				return nil, fmt.Errorf("field: action %q requires From and To", o.Action)
			}
		case "delete":
			if o.From == "" {
				return nil, fmt.Errorf("field: action %q requires From", o.Action)
			}
		default:
			return nil, fmt.Errorf("field: unknown action %q", o.Action)
		}
	}
	return &Processor{ops: ops}, nil
}

// Apply returns a new entry with all ops applied. The original is not mutated.
func (p *Processor) Apply(e reader.LogEntry) reader.LogEntry {
	extra := make(map[string]any, len(e.Extra))
	for k, v := range e.Extra {
		extra[k] = v
	}
	for _, o := range p.ops {
		switch o.Action {
		case "copy":
			if v, ok := extra[o.From]; ok {
				extra[o.To] = v
			}
		case "rename":
			if v, ok := extra[o.From]; ok {
				extra[o.To] = v
				delete(extra, o.From)
			}
		case "delete":
			delete(extra, o.From)
		}
	}
	e.Extra = extra
	return e
}
