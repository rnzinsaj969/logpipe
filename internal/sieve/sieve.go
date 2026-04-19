package sieve

import (
	"fmt"
	"regexp"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines a single sieve condition. All non-empty fields must match.
type Rule struct {
	Level   string
	Service string
	Pattern string

	re *regexp.Regexp
}

// Sieve filters log entries by applying a set of allow-rules.
// An entry passes if it matches ANY rule. If no rules are defined,
// all entries pass.
type Sieve struct {
	rules []Rule
}

// New creates a Sieve from the provided rules, compiling any regex patterns.
func New(rules []Rule) (*Sieve, error) {
	compiled := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("sieve: invalid pattern %q: %w", r.Pattern, err)
			}
			r.re = re
		}
		compiled[i] = r
	}
	return &Sieve{rules: compiled}, nil
}

// Apply returns true if the entry should be kept.
func (s *Sieve) Apply(e reader.LogEntry) bool {
	if len(s.rules) == 0 {
		return true
	}
	for _, r := range s.rules {
		if matches(r, e) {
			return true
		}
	}
	return false
}

func matches(r Rule, e reader.LogEntry) bool {
	if r.Level != "" && r.Level != e.Level {
		return false
	}
	if r.Service != "" && r.Service != e.Service {
		return false
	}
	if r.re != nil && !r.re.MatchString(e.Message) {
		return false
	}
	return true
}
