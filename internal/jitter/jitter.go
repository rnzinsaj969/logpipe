// Package jitter applies randomised time-based jitter to log entry timestamps.
// It is useful for spreading bursts of entries that share an identical timestamp
// so that downstream consumers can distinguish ordering.
package jitter

import (
	"fmt"
	"math/rand"
	"time"

	"logpipe/internal/reader"
)

// Source is the interface satisfied by rand.Rand and any test double.
type Source interface {
	Int63n(n int64) int64
}

// Jitterer adds a small random offset to each entry's timestamp.
type Jitterer struct {
	max time.Duration
	src Source
}

// New returns a Jitterer that spreads timestamps by up to max duration.
// max must be positive.
func New(max time.Duration, src Source) (*Jitterer, error) {
	if max <= 0 {
		return nil, fmt.Errorf("jitter: max must be positive, got %s", max)
	}
	if src == nil {
		src = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	}
	return &Jitterer{max: max, src: src}, nil
}

// Apply returns a copy of e with a random offset in [0, max) added to its
// timestamp. Entries with a zero timestamp are left unchanged.
func (j *Jitterer) Apply(e reader.LogEntry) reader.LogEntry {
	if e.Timestamp.IsZero() {
		return e
	}
	offset := time.Duration(j.src.Int63n(int64(j.max)))
	out := e
	out.Timestamp = e.Timestamp.Add(offset)
	return out
}
