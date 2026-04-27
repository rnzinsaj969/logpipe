// Package surge detects sudden spikes in log volume for a given service
// within a rolling time window and reports whether the current rate
// exceeds a configured multiplier of the baseline average.
package surge

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Detector tracks per-service event counts and signals when the recent
// rate surges above a threshold multiple of the long-term baseline.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	multiple  float64
	buckets   map[string][]time.Time
	clock     func() time.Time
}

// New returns a Detector that considers any burst within window that is
// more than multiple times the baseline rate a surge.
// multiple must be > 1.0 and window must be positive.
func New(window time.Duration, multiple float64) (*Detector, error) {
	return newWithClock(window, multiple, time.Now)
}

func newWithClock(window time.Duration, multiple float64, clock func() time.Time) (*Detector, error) {
	if window <= 0 {
		return nil, errors.New("surge: window must be positive")
	}
	if multiple <= 1.0 {
		return nil, errors.New("surge: multiple must be greater than 1.0")
	}
	return &Detector{
		window:   window,
		multiple: multiple,
		buckets:  make(map[string][]time.Time),
		clock:    clock,
	}, nil
}

// Record adds the entry to the detector and returns true if the service's
// recent rate constitutes a surge.
func (d *Detector) Record(e reader.LogEntry) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	svc := e.Service
	cutoff := now.Add(-d.window)

	times := d.buckets[svc]
	// evict stale entries
	start := 0
	for start < len(times) && times[start].Before(cutoff) {
		start++
	}
	times = append(times[start:], now)
	d.buckets[svc] = times

	if len(times) < 2 {
		return false
	}

	// baseline: average rate over the full window
	span := times[len(times)-1].Sub(times[0]).Seconds()
	if span <= 0 {
		return false
	}
	baseline := float64(len(times)-1) / span

	// recent rate: last quarter of the window
	recentCutoff := now.Add(-d.window / 4)
	recentCount := 0
	for _, t := range times {
		if !t.Before(recentCutoff) {
			recentCount++
		}
	}
	recentRate := float64(recentCount) / (d.window.Seconds() / 4)

	return recentRate >= baseline*d.multiple
}

// Reset clears all recorded data for every service.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buckets = make(map[string][]time.Time)
}
