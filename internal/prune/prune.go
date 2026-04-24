// Package prune removes specified keys from the Extra field of a log entry.
package prune

import (
	"errors"

	"logpipe/internal/reader"
)

// Pruner removes a fixed set of keys from entry.Extra.
type Pruner struct {
	keys map[string]struct{}
}

// New returns a Pruner that will remove each of the given keys from
// entry.Extra. At least one key must be provided.
func New(keys ...string) (*Pruner, error) {
	if len(keys) == 0 {
		return nil, errors.New("prune: at least one key is required")
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if k == "" {
			return nil, errors.New("prune: key must not be empty")
		}
		m[k] = struct{}{}
	}
	return &Pruner{keys: m}, nil
}

// Apply returns a copy of e with the configured keys removed from Extra.
// If Extra is nil or none of the keys are present the entry is returned
// unchanged.
func (p *Pruner) Apply(e reader.LogEntry) reader.LogEntry {
	if len(e.Extra) == 0 {
		return e
	}

	// Only allocate a new map when we actually find a key to remove.
	var out map[string]any
	for k, v := range e.Extra {
		if _, drop := p.keys[k]; drop {
			if out == nil {
				out = make(map[string]any, len(e.Extra))
				for k2, v2 := range e.Extra {
					out[k2] = v2
				}
			}
			delete(out, k)
		} else {
			_ = v // suppress unused warning
		}
	}
	if out == nil {
		return e
	}
	e.Extra = out
	return e
}

// Keys returns the set of keys this Pruner will remove.
func (p *Pruner) Keys() []string {
	out := make([]string, 0, len(p.keys))
	for k := range p.keys {
		out = append(out, k)
	}
	return out
}
