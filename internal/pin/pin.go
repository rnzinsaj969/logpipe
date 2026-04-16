// Package pin provides field pinning: promoting nested or extra fields to
// top-level LogEntry fields when a matching key is found.
package pin

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule describes a single pin operation.
type Rule struct {
	// Key is the extra-field key to promote.
	Key string
	// Target is one of "message", "level", or "service".
	Target string
}

// Pinner promotes extra fields to top-level entry fields.
type Pinner struct {
	rules []Rule
}

// New creates a Pinner from the supplied rules.
// Returns an error if any rule has an empty Key or unsupported Target.
func New(rules []Rule) (*Pinner, error) {
	for _, r := range rules {
		if r.Key == "" {
			return nil, errors.New("pin: rule key must not be empty")
		}
		switch r.Target {
		case "message", "level", "service":
		default:
			return nil, errors.New("pin: unsupported target " + r.Target)
		}
	}
	return &Pinner{rules: rules}, nil
}

// Apply returns a new LogEntry with matching extra fields promoted.
// The original entry is never mutated.
func (p *Pinner) Apply(e reader.LogEntry) reader.LogEntry {
	out := e
	newExtra := make(map[string]any, len(e.Extra))
	for k, v := range e.Extra {
		newExtra[k] = v
	}
	out.Extra = newExtra

	for _, r := range p.rules {
		val, ok := out.Extra[r.Key]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		switch r.Target {
		case "message":
			out.Message = s
		case "level":
			out.Level = s
		case "service":
			out.Service = s
		}
		delete(out.Extra, r.Key)
	}
	return out
}
