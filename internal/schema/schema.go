// Package schema provides log entry field validation and normalization
// for structured log entries flowing through the pipeline.
package schema

import (
	"errors"
	"strings"
	"time"
)

// Required is the set of field names that must be present in a log entry.
var Required = []string{"message", "level", "service"}

// Validator checks log entry maps against a required field schema.
type Validator struct {
	required []string
}

// New returns a Validator that enforces the given required field names.
func New(required []string) *Validator {
	r := make([]string, len(required))
	copy(r, required)
	return &Validator{required: r}
}

// Validate returns an error if any required field is missing or empty in entry.
func (v *Validator) Validate(entry map[string]any) error {
	var missing []string
	for _, field := range v.required {
		val, ok := entry[field]
		if !ok {
			missing = append(missing, field)
			continue
		}
		s, isStr := val.(string)
		if isStr && strings.TrimSpace(s) == "" {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return errors.New("missing required fields: " + strings.Join(missing, ", "))
	}
	return nil
}

// Normalize fills in default values for optional fields if absent.
// It sets "timestamp" to the current UTC time in RFC3339 format when missing.
func Normalize(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	if _, ok := out["timestamp"]; !ok {
		out["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}
	if level, ok := out["level"].(string); ok {
		out["level"] = strings.ToLower(strings.TrimSpace(level))
	}
	return out
}
