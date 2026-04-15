package label

import "github.com/logpipe/logpipe/internal/reader"

// Rule defines a label to attach to log entries whose fields match a condition.
type Rule struct {
	// Field is the log entry field to inspect ("message", "level", "service", or an extra key).
	Field string
	// Contains is a substring that must appear in the field value for the label to apply.
	Contains string
	// Label is the key added to the entry's Extra map when the rule matches.
	Label string
	// Value is the value assigned to the label key.
	Value string
}

// Labeler attaches metadata labels to log entries based on configurable rules.
type Labeler struct {
	rules []Rule
}

// New creates a Labeler with the given set of rules.
func New(rules []Rule) *Labeler {
	return &Labeler{rules: rules}
}

// Apply returns a copy of entry with any matching labels added to Extra.
// The original entry is never mutated.
func (l *Labeler) Apply(e reader.LogEntry) reader.LogEntry {
	out := e
	for _, r := range l.rules {
		if matches(e, r) {
			out = addLabel(out, r.Label, r.Value)
		}
	}
	return out
}

// HasRules reports whether the labeler has at least one rule configured.
func (l *Labeler) HasRules() bool {
	return len(l.rules) > 0
}

func matches(e reader.LogEntry, r Rule) bool {
	var val string
	switch r.Field {
	case "message":
		val = e.Message
	case "level":
		val = e.Level
	case "service":
		val = e.Service
	default:
		if e.Extra != nil {
			v, ok := e.Extra[r.Field]
			if !ok {
				return false
			}
			if s, ok := v.(string); ok {
				val = s
			}
		}
	}
	return r.Contains != "" && contains(val, r.Contains)
}

func addLabel(e reader.LogEntry, key, value string) reader.LogEntry {
	extra := make(map[string]any, len(e.Extra)+1)
	for k, v := range e.Extra {
		extra[k] = v
	}
	extra[key] = value
	e.Extra = extra
	return e
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
