// Package transform provides field-level transformations for log entries
// before they are passed downstream in the pipeline.
package transform

import (
	"strings"

	"github.com/yourorg/logpipe/internal/reader"
)

// Rule describes a single transformation to apply to a log entry field.
type Rule struct {
	// Field is the top-level field name to transform (e.g. "message", "service").
	Field string
	// Op is the operation to apply: "upper", "lower", "trim", or "prefix".
	Op string
	// Value is an optional argument used by ops such as "prefix".
	Value string
}

// Transformer applies a sequence of Rules to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that will apply the given rules in order.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply returns a copy of entry with all rules applied.
func (t *Transformer) Apply(entry reader.LogEntry) reader.LogEntry {
	for _, r := range t.rules {
		entry = applyRule(entry, r)
	}
	return entry
}

// HasRules reports whether any rules are configured.
func (t *Transformer) HasRules() bool {
	return len(t.rules) > 0
}

func applyRule(entry reader.LogEntry, r Rule) reader.LogEntry {
	switch r.Field {
	case "message":
		entry.Message = applyOp(entry.Message, r)
	case "service":
		entry.Service = applyOp(entry.Service, r)
	case "level":
		entry.Level = applyOp(entry.Level, r)
	}
	return entry
}

func applyOp(s string, r Rule) string {
	switch r.Op {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "trim":
		return strings.TrimSpace(s)
	case "prefix":
		return r.Value + s
	default:
		return s
	}
}
