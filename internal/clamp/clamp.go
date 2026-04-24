// Package clamp provides a per-service rate limiter that enforces a maximum
// number of log entries within a rolling time window, dropping entries that
// exceed the cap.
package clamp

import (
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

type bucket struct {
	count int
	reset time.Time
}

// Limiter enforces a per-service entry cap within a rolling window.
type Limiter struct {
	mu     sync.Mutex
	max    int
	window time.Duration
	clock  func() time.Time
	state  map[string]*bucket
}

// New returns a Limiter that allows at most max entries per service within
// the given window duration. It returns an error if max is not positive or
// window is zero.
func New(max int, window time.Duration) (*Limiter, error) {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, clock func() time.Time) (*Limiter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("clamp: max must be positive, got %d", max)
	}
	if window <= 0 {
		return nil, fmt.Errorf("clamp: window must be positive, got %s", window)
	}
	return &Limiter{
		max:    max,
		window: window,
		clock:  clock,
		state:  make(map[string]*bucket),
	}, nil
}

// Allow reports whether the entry should be forwarded. Entries from services
// that have exceeded the cap for the current window are dropped.
func (l *Limiter) Allow(e reader.LogEntry) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	svc := e.Service

	b, ok := l.state[svc]
	if !ok || now.After(b.reset) {
		l.state[svc] = &bucket{count: 1, reset: now.Add(l.window)}
		return true
	}

	if b.count >= l.max {
		return false
	}

	b.count++
	return true
}
