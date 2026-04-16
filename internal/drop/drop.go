package drop

import (
	"fmt"
	"regexp"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule describes a single drop condition. An entry matching any rule is discarded.
type Rule struct {
	Level   string
	Service string
	Pattern string // regexp applied to Message

	re *regexp.Regexp
}

// Dropper discards log entries that match at least one configured rule.
type Dropper struct {
	rules []Rule
}

// New compiles all pattern rules and returns a Dropper.
// Returns an error if any Pattern is an invalid regular expression.
func New(rules []Rule) (*Dropper, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		compiled[i] = r
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("drop: invalid pattern %q: %w", r.Pattern, err)
			}
			compiled[i].re = re
		}
	}
	return &Dropper{rules: compiled}, nil
}

// ShouldDrop returns true when the entry matches at least one rule.
func (d *Dropper) ShouldDrop(e reader.LogEntry) bool {
	for _, r := range d.rules {
		if r.Level != "" && r.Level != e.Level {
			continue
		}
		if r.Service != "" && r.Service != e.Service {
			continue
		}
		if r.re != nil && !r.re.MatchString(e.Message) {
			continue
		}
		return true
	}
	return false
}

// Apply returns a filtered slice that excludes dropped entries.
func (d *Dropper) Apply(entries []reader.LogEntry) []reader.LogEntry {
	out := entries[:0:0]
	for _, e := range entries {
		if !d.ShouldDrop(e) {
			out = append(out, e)
		}
	}
	return out
}
