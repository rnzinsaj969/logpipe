// Package clip provides field-level value clipping for log entries,
// clamping numeric fields to a specified [min, max] range.
package clip

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rule defines a clipping rule for a single extra field.
type Rule struct {
	// Field is the key inside LogEntry.Extra to clamp.
	Field string
	// Min is the lower bound (inclusive).
	Min float64
	// Max is the upper bound (inclusive).
	Max float64
}

// Clipper applies numeric clamping rules to log entry extra fields.
type Clipper struct {
	rules []Rule
}

// New creates a Clipper from the provided rules.
// It returns an error if any rule has Min > Max or an empty Field.
func New(rules []Rule) (*Clipper, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("clip: rule %d has empty field", i)
		}
		if r.Min > r.Max {
			return nil, fmt.Errorf("clip: rule %d has min %.4g > max %.4g", i, r.Min, r.Max)
		}
	}
	return &Clipper{rules: rules}, nil
}

// Rules returns a copy of the rules configured on this Clipper.
func (c *Clipper) Rules() []Rule {
	out := make([]Rule, len(c.rules))
	copy(out, c.rules)
	return out
}

// Apply returns a new LogEntry with numeric extra fields clamped according
// to the configured rules. Non-numeric values are left untouched.
func (c *Clipper) Apply(e reader.LogEntry) reader.LogEntry {
	if len(c.rules) == 0 || len(e.Extra) == 0 {
		return e
	}

	newExtra := make(map[string]any, len(e.Extra))
	for k, v := range e.Extra {
		newExtra[k] = v
	}

	for _, r := range c.rules {
		v, ok := newExtra[r.Field]
		if !ok {
			continue
		}
		f, ok := toFloat(v)
		if !ok {
			continue
		}
		newExtra[r.Field] = clamp(f, r.Min, r.Max)
	}

	e.Extra = newExtra
	return e
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
