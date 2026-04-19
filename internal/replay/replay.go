package replay

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Replayer holds a fixed set of log entries and replays them at a
// configurable speed multiplier relative to their original timestamps.
type Replayer struct {
	mu      sync.Mutex
	entries []reader.LogEntry
	speed   float64
}

// Options configures a Replayer.
type Options struct {
	// Speed is the playback multiplier. 1.0 = real-time, 2.0 = double speed.
	// Values <= 0 default to 1.0.
	Speed float64
}

// New creates a Replayer from a slice of entries.
func New(entries []reader.LogEntry, opts Options) (*Replayer, error) {
	if len(entries) == 0 {
		return nil, errors.New("replay: entries must not be empty")
	}
	speed := opts.Speed
	if speed <= 0 {
		speed = 1.0
	}
	cp := make([]reader.LogEntry, len(entries))
	copy(cp, entries)
	return &Replayer{entries: cp, speed: speed}, nil
}

// Run sends entries to out in timestamp order, sleeping between them
// according to the speed multiplier. It stops when ctx is done.
func (r *Replayer) Run(ctx interface{ Done() <-chan struct{} }, out chan<- reader.LogEntry) {
	r.mu.Lock()
	entries := r.entries
	speed := r.speed
	r.mu.Unlock()

	for i, e := range entries {
		if i > 0 {
			prev := entries[i-1].Timestamp
			gap := e.Timestamp.Sub(prev)
			if gap > 0 {
				delay := time.Duration(float64(gap) / speed)
				select {
				case <-ctx.Done():
					return
				case <-time.After(delay):
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		case out <- e:
		}
	}
}

// Len returns the number of stored entries.
func (r *Replayer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}
