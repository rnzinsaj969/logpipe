package reorder

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		WindowSize: 5,
		MaxAge:     2 * time.Second,
	}
}

// Options controls reorder behaviour.
type Options struct {
	// WindowSize is the maximum number of entries buffered before flushing.
	WindowSize int
	// MaxAge is the maximum time an entry may be held before it is flushed.
	MaxAge time.Duration
}

// Reorderer buffers log entries and emits them sorted by timestamp.
type Reorderer struct {
	mu      sync.Mutex
	opts    Options
	buf     []reader.LogEntry
	clock   func() time.Time
}

// New returns a Reorderer with the given options.
func New(opts Options) (*Reorderer, error) {
	if opts.WindowSize <= 0 {
		return nil, errors.New("reorder: WindowSize must be greater than zero")
	}
	if opts.MaxAge <= 0 {
		return nil, errors.New("reorder: MaxAge must be greater than zero")
	}
	return &Reorderer{
		opts:  opts,
		clock: time.Now,
	}, nil
}

func newWithClock(opts Options, clock func() time.Time) (*Reorderer, error) {
	r, err := New(opts)
	if err != nil {
		return nil, err
	}
	r.clock = clock
	return r, nil
}

// Add buffers an entry. It returns any entries that are ready to be emitted.
func (r *Reorderer) Add(entry reader.LogEntry) []reader.LogEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf = append(r.buf, entry)
	if len(r.buf) >= r.opts.WindowSize {
		return r.flush()
	}
	return nil
}

// Flush forces all buffered entries to be emitted in timestamp order.
func (r *Reorderer) Flush() []reader.LogEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.flush()
}

// Drain emits entries whose age exceeds MaxAge, sorted by timestamp.
func (r *Reorderer) Drain() []reader.LogEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.clock()
	var keep, emit []reader.LogEntry
	for _, e := range r.buf {
		if !e.Timestamp.IsZero() && now.Sub(e.Timestamp) >= r.opts.MaxAge {
			emit = append(emit, e)
		} else {
			keep = append(keep, e)
		}
	}
	r.buf = keep
	sortEntries(emit)
	return emit
}

func (r *Reorderer) flush() []reader.LogEntry {
	out := make([]reader.LogEntry, len(r.buf))
	copy(out, r.buf)
	r.buf = r.buf[:0]
	sortEntries(out)
	return out
}

func sortEntries(entries []reader.LogEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
}
