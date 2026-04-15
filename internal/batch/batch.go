package batch

import (
	"errors"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// ErrEmptyBatch is returned when Flush is called with no entries.
var ErrEmptyBatch = errors.New("batch: no entries to flush")

// Options configures the Batcher.
type Options struct {
	// MaxSize is the maximum number of entries before an automatic flush.
	MaxSize int
	// MaxAge is the maximum duration entries are held before an automatic flush.
	MaxAge time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxSize: 100,
		MaxAge:  5 * time.Second,
	}
}

// Batcher accumulates log entries and flushes them as a slice when either
// MaxSize is reached or MaxAge elapses since the first entry was added.
type Batcher struct {
	opts    Options
	buf     []reader.LogEntry
	first   time.Time
	clock   func() time.Time
}

// New creates a Batcher with the given options.
func New(opts Options) *Batcher {
	if opts.MaxSize <= 0 {
		opts.MaxSize = DefaultOptions().MaxSize
	}
	if opts.MaxAge <= 0 {
		opts.MaxAge = DefaultOptions().MaxAge
	}
	return &Batcher{opts: opts, clock: time.Now}
}

// Add appends an entry to the current batch.
// It reports whether the batch is ready to flush.
func (b *Batcher) Add(e reader.LogEntry) bool {
	if len(b.buf) == 0 {
		b.first = b.clock()
	}
	b.buf = append(b.buf, e)
	return b.ready()
}

// Ready reports whether the batch should be flushed.
func (b *Batcher) Ready() bool { return b.ready() }

func (b *Batcher) ready() bool {
	if len(b.buf) == 0 {
		return false
	}
	if len(b.buf) >= b.opts.MaxSize {
		return true
	}
	return b.clock().Sub(b.first) >= b.opts.MaxAge
}

// Flush returns the accumulated entries and resets the batch.
// Returns ErrEmptyBatch if there are no entries.
func (b *Batcher) Flush() ([]reader.LogEntry, error) {
	if len(b.buf) == 0 {
		return nil, ErrEmptyBatch
	}
	out := make([]reader.LogEntry, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	b.first = time.Time{}
	return out, nil
}

// Len returns the number of buffered entries.
func (b *Batcher) Len() int { return len(b.buf) }
