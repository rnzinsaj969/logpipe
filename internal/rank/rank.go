// Package rank assigns a numeric priority to log entries based on
// configurable field-value rules. Higher scores indicate higher priority.
package rank

import (
	"errors"
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule maps a field/value condition to a priority score.
type Rule struct {
	// Field is one of "level", "service", or an extra field key.
	Field string
	// Value is the expected value (case-insensitive).
	Value string
	// Score is added to the entry's total when the rule matches.
	Score int
}

// Ranker evaluates rules against log entries and returns a priority score.
type Ranker struct {
	rules []Rule
	field string
}

// New creates a Ranker that writes the computed score into the given extra
// field name. It returns an error if rules is empty or field is blank.
func New(field string, rules []Rule) (*Ranker, error) {
	if strings.TrimSpace(field) == "" {
		return nil, errors.New("rank: output field must not be empty")
	}
	if len(rules) == 0 {
		return nil, errors.New("rank: at least one rule is required")
	}
	return &Ranker{rules: rules, field: field}, nil
}

// Apply evaluates all rules against e and stores the total score in the
// configured extra field. The original entry is never mutated.
func (r *Ranker) Apply(e reader.LogEntry) reader.LogEntry {
	total := 0
	for _, rule := range r.rules {
		if r.matches(e, rule) {
			total += rule.Score
		}
	}

	extra := make(map[string]any, len(e.Extra)+1)
	for k, v := range e.Extra {
		extra[k] = v
	}
	extra[r.field] = total
	return reader.LogEntry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Service:   e.Service,
		Message:   e.Message,
		Extra:     extra,
	}
}

func (r *Ranker) matches(e reader.LogEntry, rule Rule) bool {
	var actual string
	switch rule.Field {
	case "level":
		actual = e.Level
	case "service":
		actual = e.Service
	default:
		if v, ok := e.Extra[rule.Field]; ok {
			actual = strings.TrimSpace(fmt.Sprintf("%v", v))
		}
	}
	return strings.EqualFold(actual, rule.Value)
}
