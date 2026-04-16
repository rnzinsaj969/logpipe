package quota

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a service has exceeded its quota.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Options configures the Quota limiter.
type Options struct {
	MaxEntries int
	Window     time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxEntries: 1000,
		Window:     time.Minute,
	}
}

type bucket struct {
	count int
	reset time.Time
}

// Quota enforces a maximum number of log entries per service per time window.
type Quota struct {
	mu      sync.Mutex
	opts    Options
	clock   func() time.Time
	buckets map[string]*bucket
}

// New creates a Quota limiter with the given options.
func New(opts Options) (*Quota, error) {
	if opts.MaxEntries <= 0 {
		return nil, errors.New("quota: MaxEntries must be positive")
	}
	if opts.Window <= 0 {
		return nil, errors.New("quota: Window must be positive")
	}
	return &Quota{
		opts:    opts,
		clock:   time.Now,
		buckets: make(map[string]*bucket),
	}, nil
}

// Allow reports whether the given service is within its quota.
// It increments the counter and returns ErrQuotaExceeded if the limit is reached.
func (q *Quota) Allow(service string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := q.clock()
	b, ok := q.buckets[service]
	if !ok || now.After(b.reset) {
		q.buckets[service] = &bucket{count: 1, reset: now.Add(q.opts.Window)}
		return nil
	}
	if b.count >= q.opts.MaxEntries {
		return ErrQuotaExceeded
	}
	b.count++
	return nil
}

// Reset clears the quota state for a service.
func (q *Quota) Reset(service string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.buckets, service)
}

// Snapshot returns a copy of current counts keyed by service.
func (q *Quota) Snapshot() map[string]int {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make(map[string]int, len(q.buckets))
	for k, b := range q.buckets {
		out[k] = b.count
	}
	return out
}
