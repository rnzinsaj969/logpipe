// Package timeslot groups log entries into fixed-duration time buckets.
package timeslot

import (
	"errors"
	"sync"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Slot holds entries that fall within a single time bucket.
type Slot struct {
	Start   time.Time
	Entries []reader.LogEntry
}

// Bucketer partitions log entries into time slots of a fixed duration.
type Bucketer struct {
	mu       sync.Mutex
	size     time.Duration
	slots    map[int64]*Slot
	clock    func() time.Time
}

// Options configures the Bucketer.
type Options struct {
	Size  time.Duration
	Clock func() time.Time
}

// New returns a Bucketer with the given slot size.
func New(opts Options) (*Bucketer, error) {
	if opts.Size <= 0 {
		return nil, errors.New("timeslot: size must be positive")
	}
	clk := opts.Clock
	if clk == nil {
		clk = time.Now
	}
	return &Bucketer{
		size:  opts.Size,
		slots: make(map[int64]*Slot),
		clock: clk,
	}, nil
}

// Add places the entry into the appropriate time slot.
func (b *Bucketer) Add(e reader.LogEntry) {
	t := e.Timestamp
	if t.IsZero() {
		t = b.clock()
	}
	key := t.Truncate(b.size).UnixNano()
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.slots[key]; !ok {
		b.slots[key] = &Slot{Start: t.Truncate(b.size)}
	}
	b.slots[key].Entries = append(b.slots[key].Entries, e)
}

// Snapshot returns all current slots and resets internal state.
func (b *Bucketer) Snapshot() []Slot {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Slot, 0, len(b.slots))
	for _, s := range b.slots {
		out = append(out, *s)
	}
	b.slots = make(map[int64]*Slot)
	return out
}
