// Package sampling provides probabilistic log entry sampling
// to reduce volume when log rates are high.
package sampling

import (
	"math/rand"
	"sync"
)

// Sampler decides whether a log entry should be kept based on a
// configured sampling rate between 0.0 (drop all) and 1.0 (keep all).
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New returns a Sampler with the given rate. Rate is clamped to [0.0, 1.0].
func New(rate float64, src rand.Source) *Sampler {
	if rate < 0.0 {
		rate = 0.0
	}
	if rate > 1.0 {
		rate = 1.0
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(src),
	}
}

// Keep returns true if the entry should be forwarded downstream.
// It is safe for concurrent use.
func (s *Sampler) Keep() bool {
	if s.rate >= 1.0 {
		return true
	}
	if s.rate <= 0.0 {
		return false
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Rate returns the current sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}

// SetRate updates the sampling rate at runtime. Value is clamped to [0.0, 1.0].
func (s *Sampler) SetRate(rate float64) {
	if rate < 0.0 {
		rate = 0.0
	}
	if rate > 1.0 {
		rate = 1.0
	}
	s.mu.Lock()
	s.rate = rate
	s.mu.Unlock()
}
