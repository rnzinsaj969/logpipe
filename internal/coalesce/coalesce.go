package coalesce

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines how to merge a set of source fields into a single target field.
type Rule struct {
	// Sources is an ordered list of Extra field keys to check.
	Sources []string
	// Target is the Extra field key to write the first non-empty value into.
	Target string
}

// Coalescer merges multiple optional fields into a single canonical field.
type Coalescer struct {
	rules []Rule
}

// New creates a Coalescer from the provided rules.
// Returns an error if any rule has an empty Target or no Sources.
func New(rules []Rule) (*Coalescer, error) {
	for i, r := range rules {
		if r.Target == "" {
			return nil, errors.New("coalesce: rule has empty target")
		}
		if len(r.Sources) == 0 {
			return nil, errors.New("coalesce: rule has no sources")
		}
		_ = i
	}
	return &Coalescer{rules: rules}, nil
}

// Apply returns a new LogEntry with coalesced fields applied.
// The original entry is never mutated.
func (c *Coalescer) Apply(entry reader.LogEntry) reader.LogEntry {
	out := entry
	newExtra := make(map[string]any, len(entry.Extra))
	for k, v := range entry.Extra {
		newExtra[k] = v
	}

	for _, rule := range c.rules {
		for _, src := range rule.Sources {
			v, ok := newExtra[src]
			if !ok {
				continue
			}
			s, ok := v.(string)
			if !ok || s == "" {
				continue
			}
			newExtra[rule.Target] = s
			break
		}
	}

	out.Extra = newExtra
	return out
}

// HasRules reports whether the Coalescer has at least one rule.
func (c *Coalescer) HasRules() bool {
	return len(c.rules) > 0
}
