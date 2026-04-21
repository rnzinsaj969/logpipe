// Package stash provides a named entry store for temporarily holding
// log entries under a string key, with optional capacity limits.
package stash

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// DefaultCapacity is the maximum number of entries held per key when no
// explicit capacity is provided.
const DefaultCapacity = 256

// Stash holds slices of log entries keyed by an arbitrary string label.
type Stash struct {
	mu       sync.Mutex
	buckets  map[string][]reader.LogEntry
	capacity int
}

// New returns a Stash with the given per-bucket capacity.
// capacity must be greater than zero.
func New(capacity int) (*Stash, error) {
	if capacity <= 0 {
		return nil, errors.New("stash: capacity must be greater than zero")
	}
	return &Stash{
		buckets:  make(map[string][]reader.LogEntry),
		capacity: capacity,
	}, nil
}

// Put stores entry under key. If the bucket is already at capacity the
// oldest entry is evicted (FIFO).
func (s *Stash) Put(key string, entry reader.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bucket := s.buckets[key]
	if len(bucket) >= s.capacity {
		bucket = bucket[1:]
	}
	s.buckets[key] = append(bucket, entry)
}

// Flush returns all entries stored under key and removes the bucket.
// If the key does not exist an empty slice is returned.
func (s *Stash) Flush(key string) []reader.LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	bucket := s.buckets[key]
	if len(bucket) == 0 {
		return []reader.LogEntry{}
	}
	out := make([]reader.LogEntry, len(bucket))
	copy(out, bucket)
	delete(s.buckets, key)
	return out
}

// Len returns the number of entries currently held under key.
func (s *Stash) Len(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.buckets[key])
}

// Keys returns a snapshot of all keys that currently hold at least one entry.
func (s *Stash) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.buckets))
	for k := range s.buckets {
		keys = append(keys, k)
	}
	return keys
}
