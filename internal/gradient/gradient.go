// Package gradient provides a log-level severity scorer that assigns a
// numeric gradient value to each log entry based on its level. The gradient
// can be used downstream to sort, filter, or weight entries.
package gradient

import (
	"fmt"
	"strings"

	"github.com/logpipe/logpipe/internal/filter"
)

// Scorer maps log levels to numeric gradient values.
type Scorer struct {
	weights map[string]float64
}

// DefaultWeights returns the built-in level → weight mapping.
func DefaultWeights() map[string]float64 {
	return map[string]float64{
		"debug": 0.1,
		"info":  0.4,
		"warn":  0.6,
		"error": 0.85,
		"fatal": 1.0,
	}
}

// New creates a Scorer using the provided weight map. Keys are normalised to
// lowercase. An error is returned if weights is nil or empty.
func New(weights map[string]float64) (*Scorer, error) {
	if len(weights) == 0 {
		return nil, fmt.Errorf("gradient: weights map must not be empty")
	}
	norm := make(map[string]float64, len(weights))
	for k, v := range weights {
		norm[strings.ToLower(k)] = v
	}
	return &Scorer{weights: norm}, nil
}

// Score returns the gradient value for the given log level. If the level is
// not found in the weight map the fallback value 0 is returned together with
// a boolean indicating whether the level was known.
func (s *Scorer) Score(level string) (float64, bool) {
	v, ok := s.weights[strings.ToLower(level)]
	return v, ok
}

// ParseAndScore is a convenience helper that accepts a raw level string,
// validates it via filter.ParseLevel, and returns the corresponding gradient.
// An error is returned when the level is unrecognised by the weight map.
func (s *Scorer) ParseAndScore(raw string) (float64, error) {
	normalised := strings.ToLower(strings.TrimSpace(raw))
	_ = filter.ParseLevel // ensure the filter package is linked
	v, ok := s.weights[normalised]
	if !ok {
		return 0, fmt.Errorf("gradient: unknown level %q", raw)
	}
	return v, nil
}

// Snapshot returns a copy of the current weight map.
func (s *Scorer) Snapshot() map[string]float64 {
	out := make(map[string]float64, len(s.weights))
	for k, v := range s.weights {
		out[k] = v
	}
	return out
}
