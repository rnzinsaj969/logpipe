// Package offset tracks the read position (byte offset) for each log source
// so that logpipe can resume reading from where it left off after a restart.
package offset

import (
	"sync"
)

// Store holds per-source byte offsets in memory.
type Store struct {
	mu      sync.RWMutex
	offsets map[string]int64
}

// New returns an empty Store.
func New() *Store {
	return &Store{offsets: make(map[string]int64)}
}

// Set records the current offset for the given source.
func (s *Store) Set(source string, offset int64) {
	if offset < 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets[source] = offset
}

// Get returns the stored offset for source, or 0 if not found.
func (s *Store) Get(source string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.offsets[source]
}

// Delete removes the offset entry for source.
func (s *Store) Delete(source string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.offsets, source)
}

// Snapshot returns a shallow copy of all tracked offsets.
func (s *Store) Snapshot() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int64, len(s.offsets))
	for k, v := range s.offsets {
		out[k] = v
	}
	return out
}

// Reset clears all stored offsets.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.offsets = make(map[string]int64)
}
