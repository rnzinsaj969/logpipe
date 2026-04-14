// Package throttle provides per-service log entry throttling based on
// configurable time windows and maximum entry counts.
package throttle

import (
	"sync"
	"time"
)

// Clock is a function that returns the current time, injectable for testing.
type Clock func() time.Time

// Throttler tracks entry counts per service within a sliding window and
// drops entries that exceed the configured limit.
type Throttler struct {
	mu       sync.Mutex
	window   time.Duration
	maxCount int
	clock    Clock
	buckets  map[string][]time.Time
}

// New creates a Throttler that allows at most maxCount entries per service
// within the given window duration.
func New(window time.Duration, maxCount int, clock Clock) *Throttler {
	if clock == nil {
		clock = time.Now
	}
	return &Throttler{
		window:   window,
		maxCount: maxCount,
		clock:    clock,
		buckets:  make(map[string][]time.Time),
	}
}

// Allow returns true if the entry for the given service is within the allowed
// rate, and false if it should be dropped. It is safe for concurrent use.
func (t *Throttler) Allow(service string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	cutoff := now.Add(-t.window)

	times := t.buckets[service]
	// evict timestamps outside the window
	valid := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= t.maxCount {
		t.buckets[service] = valid
		return false
	}

	t.buckets[service] = append(valid, now)
	return true
}

// Reset clears all tracked state for every service.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = make(map[string][]time.Time)
}
