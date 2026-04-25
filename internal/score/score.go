// Package score assigns a numeric priority score to log entries based on
// configurable field weights. Higher scores indicate higher priority.
package score

import (
	"errors"
	"strings"

	"logpipe/internal/reader"
)

// Rule maps a field value pattern to a numeric weight contribution.
type Rule struct {
	// Field is one of "level", "service", or an extra field key.
	Field string
	// Value is the exact value to match (case-insensitive).
	Value string
	// Weight is added to the entry score when the rule matches.
	Weight float64
}

// Scorer computes a priority score for each log entry.
type Scorer struct {
	rules []Rule
}

// New creates a Scorer from the provided rules.
// Returns an error if rules is empty or any rule has a blank Field.
func New(rules []Rule) (*Scorer, error) {
	if len(rules) == 0 {
		return nil, errors.New("score: at least one rule is required")
	}
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("score: rule %d has empty field", i)
		}
	}
	return &Scorer{rules: rules}, nil
}

// Apply computes the cumulative score for e and stores it in
// e.Extra["_score"] (as a float64), then returns the modified copy.
func (s *Scorer) Apply(e reader.LogEntry) reader.LogEntry {
	var total float64
	for _, r := range s.rules {
		if s.fieldValue(e, r.Field) == strings.ToLower(r.Value) {
			total += r.Weight
		}
	}
	out := e
	if out.Extra == nil {
		out.Extra = make(map[string]any)
	} else {
		copy := make(map[string]any, len(e.Extra)+1)
		for k, v := range e.Extra {
			copy[k] = v
		}
		out.Extra = copy
	}
	out.Extra["_score"] = total
	return out
}

// Score returns the numeric score for e without mutating it.
func (s *Scorer) Score(e reader.LogEntry) float64 {
	var total float64
	for _, r := range s.rules {
		if s.fieldValue(e, r.Field) == strings.ToLower(r.Value) {
			total += r.Weight
		}
	}
	return total
}

func (s *Scorer) fieldValue(e reader.LogEntry, field string) string {
	switch field {
	case "level":
		return strings.ToLower(e.Level)
	case "service":
		return strings.ToLower(e.Service)
	default:
		if e.Extra == nil {
			return ""
		}
		v, _ := e.Extra[field]
		s, _ := v.(string)
		return strings.ToLower(s)
	}
}
