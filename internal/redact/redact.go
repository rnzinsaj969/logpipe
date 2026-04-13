package redact

import (
	"regexp"
	"strings"
)

// Rule defines a single redaction pattern and its replacement.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of redaction rules to log message strings.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor from the provided rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// NewFromPatterns compiles string patterns into a Redactor.
// Each pattern is paired with the given replacement token.
func NewFromPatterns(patterns []string, replacement string) (*Redactor, error) {
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Pattern: re, Replacement: replacement})
	}
	return New(rules), nil
}

// Apply returns a copy of s with all rule patterns replaced.
func (r *Redactor) Apply(s string) string {
	for _, rule := range r.rules {
		s = rule.Pattern.ReplaceAllString(s, rule.Replacement)
	}
	return s
}

// ApplyMap redacts values in a string map in-place, returning a new map.
func (r *Redactor) ApplyMap(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = r.Apply(v)
	}
	return out
}

// HasRules reports whether the Redactor has any active rules.
func (r *Redactor) HasRules() bool {
	return len(r.rules) > 0
}

// DefaultPatterns returns commonly used sensitive-data patterns.
func DefaultPatterns() []string {
	return []string{
		`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`,
		`(?i)password=[^\s&]+`,
		`\b(?:\d[ -]?){13,16}\b`,
		`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`,
	}
}

// DefaultReplacement is the token used when no replacement is specified.
const DefaultReplacement = "[REDACTED]"

// NewDefault builds a Redactor using DefaultPatterns and DefaultReplacement.
func NewDefault() (*Redactor, error) {
	return NewFromPatterns(DefaultPatterns(), DefaultReplacement)
}

// RedactString is a convenience wrapper for one-off redaction.
func RedactString(s string) string {
	r, err := NewDefault()
	if err != nil {
		return s
	}
	return r.Apply(s)
}

// Mask replaces the middle portion of s with asterisks, preserving
// the first and last n characters when len(s) > 2*n.
func Mask(s string, n int) string {
	if len(s) <= 2*n {
		return strings.Repeat("*", len(s))
	}
	return s[:n] + strings.Repeat("*", len(s)-2*n) + s[len(s)-n:]
}
