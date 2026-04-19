// Package suppress provides a processor that silences repeated identical
// log entries within a configurable cooldown window per service.
package suppress

import (
	"errors"
	"sync"
	"time"

	"logpipe/internal/reader"
)

// Options configures the Suppressor.
type Options struct {
	// Cooldown is the minimum duration between identical messages per service.
	Cooldown time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Cooldown: 5 * time.Second}
}

type key struct {
	service string
	message string
}

// Suppressor drops repeated log entries within the cooldown window.
type Suppressor struct {
	opts  Options
	clock func() time.Time
	mu    sync.Mutex
	seen  map[key]time.Time
}

// New creates a Suppressor with the given options.
func New(opts Options) (*Suppressor, error) {
	if opts.Cooldown <= 0 {
		return nil, errors.New("suppress: cooldown must be positive")
	}
	return &Suppressor{
		opts:  opts,
		clock: time.Now,
		seen:  make(map[key]time.Time),
	}, nil
}

// Apply returns false when the entry should be suppressed, true otherwise.
func (s *Suppressor) Apply(e reader.LogEntry) bool {
	k := key{service: e.Service, message: e.Message}
	now := s.clock()
	s.mu.Lock()
	defer s.mu.Unlock()
	if last, ok := s.seen[k]; ok && now.Sub(last) < s.opts.Cooldown {
		return false
	}
	s.seen[k] = now
	return true
}

// Evict removes stale entries older than the cooldown to bound memory.
func (s *Suppressor) Evict() {
	now := s.clock()
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, t := range s.seen {
		if now.Sub(t) >= s.opts.Cooldown {
			delete(s.seen, k)
		}
	}
}
