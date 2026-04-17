// Package journal provides an append-only, in-memory event journal for
// recording log entries with sequence numbers and optional source tagging.
package journal

import (
	"errors"
	"sync"

	"github.com/logpipe/internal/reader"
)

// Entry wraps a log entry with a monotonic sequence number and source tag.
type Entry struct {
	Seq    uint64
	Source string
	Log    reader.LogEntry
}

// Journal is a bounded, append-only in-memory store of log entries.
type Journal struct {
	mu      sync.RWMutex
	entries []Entry
	cap     int
	seq     uint64
}

// New creates a Journal with the given maximum capacity.
// capacity must be greater than zero.
func New(capacity int) (*Journal, error) {
	if capacity <= 0 {
		return nil, errors.New("journal: capacity must be greater than zero")
	}
	return &Journal{cap: capacity}, nil
}

// Append adds a log entry to the journal under the given source tag.
// When the journal is full the oldest entry is evicted.
func (j *Journal) Append(source string, log reader.LogEntry) Entry {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.seq++
	e := Entry{Seq: j.seq, Source: source, Log: log}
	if len(j.entries) >= j.cap {
		j.entries = j.entries[1:]
	}
	j.entries = append(j.entries, e)
	return e
}

// Snapshot returns a copy of all current entries in insertion order.
func (j *Journal) Snapshot() []Entry {
	j.mu.RLock()
	defer j.mu.RUnlock()
	out := make([]Entry, len(j.entries))
	copy(out, j.entries)
	return out
}

// Len returns the number of entries currently held.
func (j *Journal) Len() int {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return len(j.entries)
}

// Clear removes all entries from the journal.
func (j *Journal) Clear() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.entries = j.entries[:0]
}
