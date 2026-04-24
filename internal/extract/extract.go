// Package extract provides a processor that promotes nested extra fields
// to top-level log entry fields (message, level, or service).
package extract

import (
	"errors"
	"fmt"

	"logpipe/internal/reader"
)

// Field names that can be targeted for promotion.
const (
	TargetMessage = "message"
	TargetLevel   = "level"
	TargetService = "service"
)

// Rule describes a single extraction: copy the value of an extra field
// into a top-level entry field.
type Rule struct {
	// From is the key inside LogEntry.Extra to read.
	From string
	// To is the top-level field to write (message, level, or service).
	To string
	// Overwrite controls whether an existing non-empty target is replaced.
	Overwrite bool
}

// Extractor applies extraction rules to log entries.
type Extractor struct {
	rules []Rule
}

// New creates an Extractor from the given rules.
// Returns an error if any rule has an empty From/To or an unsupported target.
func New(rules []Rule) (*Extractor, error) {
	for i, r := range rules {
		if r.From == "" {
			return nil, fmt.Errorf("extract: rule %d: From must not be empty", i)
		}
		if r.To != TargetMessage && r.To != TargetLevel && r.To != TargetService {
			return nil, fmt.Errorf("extract: rule %d: unsupported target %q", i, r.To)
		}
	}
	if len(rules) == 0 {
		return nil, errors.New("extract: at least one rule is required")
	}
	return &Extractor{rules: rules}, nil
}

// Apply returns a new LogEntry with extraction rules applied.
// The original entry is never mutated.
func (e *Extractor) Apply(entry reader.LogEntry) reader.LogEntry {
	out := entry
	out.Extra = make(map[string]any, len(entry.Extra))
	for k, v := range entry.Extra {
		out.Extra[k] = v
	}

	for _, r := range e.rules {
		val, ok := out.Extra[r.From]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		switch r.To {
		case TargetMessage:
			if out.Message == "" || r.Overwrite {
				out.Message = s
			}
		case TargetLevel:
			if out.Level == "" || r.Overwrite {
				out.Level = s
			}
		case TargetService:
			if out.Service == "" || r.Overwrite {
				out.Service = s
			}
		}
	}
	return out
}
