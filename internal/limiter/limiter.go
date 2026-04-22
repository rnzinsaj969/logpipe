// Package limiter provides a per-service concurrent entry limiter that caps
// the number of log entries being processed simultaneously for any given
// service. Entries that exceed the configured concurrency limit are dropped.
package limiter

import (
	"errors"
	"sync"

	"github.com/logpipe/internal/reader"
)

// Limiter tracks in-flight entry counts per service and enforces a maximum
// concurrency ceiling.
type Limiter struct {
	mu      sync.Mutex
	max     int
	counts  map[string]int
}

// New creates a Limiter that allows at most max concurrent entries per service.
// max must be greater than zero.
func New(max int) (*Limiter, error) {
	if max <= 0 {
		return nil, errors.New("limiter: max must be greater than zero")
	}
	return &Limiter{
		max:    max,
		counts: make(map[string]int),
	}, nil
}

// Acquire attempts to acquire a slot for the service identified by entry.
// It returns true when a slot is available, false when the limit is reached.
func (l *Limiter) Acquire(entry reader.LogEntry) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	svc := entry.Service
	if l.counts[svc] >= l.max {
		return false
	}
	l.counts[svc]++
	return true
}

// Release decrements the in-flight count for the service identified by entry.
// It is a no-op if the count is already zero.
func (l *Limiter) Release(entry reader.LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	svc := entry.Service
	if l.counts[svc] > 0 {
		l.counts[svc]--
	}
}

// Snapshot returns a copy of the current in-flight counts keyed by service.
func (l *Limiter) Snapshot() map[string]int {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[string]int, len(l.counts))
	for k, v := range l.counts {
		out[k] = v
	}
	return out
}
