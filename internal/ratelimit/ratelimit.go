// Package ratelimit provides a token-bucket rate limiter for controlling
// log entry throughput across pipeline sources.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which log entries are processed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to maxPerSec entries per second.
// It uses a token-bucket algorithm with a burst capacity equal to maxPerSec.
func New(maxPerSec float64) *Limiter {
	now := time.Now()
	return &Limiter{
		tokens:   maxPerSec,
		max:      maxPerSec,
		rate:     maxPerSec,
		lastTick: now,
		clock:    time.Now,
	}
}

// Allow reports whether one token can be consumed.
// It refills tokens proportionally to elapsed time since the last call.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the limiter to a full token bucket.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.lastTick = l.clock()
}
