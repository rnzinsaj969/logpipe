package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds retry policy parameters.
type Config struct {
	// MaxAttempts is the total number of tries (including the first).
	MaxAttempts int
	// BaseDelay is the initial back-off duration.
	BaseDelay time.Duration
	// MaxDelay caps the exponential back-off.
	MaxDelay time.Duration
}

// Retryer executes a function with exponential back-off.
type Retryer struct {
	cfg   Config
	sleep func(time.Duration) // injectable for testing
}

// New returns a Retryer with the given Config.
func New(cfg Config) *Retryer {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if cfg.BaseDelay <= 0 {
		cfg.BaseDelay = 100 * time.Millisecond
	}
	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = 5 * time.Second
	}
	return &Retryer{
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

// Do calls fn up to MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < r.cfg.MaxAttempts-1 {
			delay := r.backoff(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}
	return errors.Join(ErrMaxAttempts, lastErr)
}

// backoff computes the delay for the given attempt index.
func (r *Retryer) backoff(attempt int) time.Duration {
	delay := float64(r.cfg.BaseDelay) * math.Pow(2, float64(attempt))
	if delay > float64(r.cfg.MaxDelay) {
		delay = float64(r.cfg.MaxDelay)
	}
	return time.Duration(delay)
}
