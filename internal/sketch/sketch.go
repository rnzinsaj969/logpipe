// Package sketch provides a probabilistic frequency estimator using a
// Count-Min Sketch, allowing approximate per-key occurrence counts over
// a stream of log entries without unbounded memory growth.
package sketch

import (
	"errors"
	"hash/fnv"
	"sync"
)

const defaultDepth = 4

// Sketch is a Count-Min Sketch that tracks approximate hit counts for
// arbitrary string keys.
type Sketch struct {
	mu    sync.Mutex
	table [][]uint64
	width uint64
	depth int
}

// New returns a Sketch with the given width (number of buckets per row).
// A larger width reduces collision probability at the cost of memory.
// width must be at least 1.
func New(width uint64) (*Sketch, error) {
	if width == 0 {
		return nil, errors.New("sketch: width must be at least 1")
	}
	table := make([][]uint64, defaultDepth)
	for i := range table {
		table[i] = make([]uint64, width)
	}
	return &Sketch{table: table, width: width, depth: defaultDepth}, nil
}

// Add increments the count for key by delta. delta must be positive;
// non-positive values are ignored.
func (s *Sketch) Add(key string, delta uint64) {
	if delta == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < s.depth; i++ {
		s.table[i][s.bucket(key, i)] += delta
	}
}

// Count returns the estimated number of times key has been added.
// The result is an upper bound on the true count.
func (s *Sketch) Count(key string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var min uint64
	for i := 0; i < s.depth; i++ {
		v := s.table[i][s.bucket(key, i)]
		if i == 0 || v < min {
			min = v
		}
	}
	return min
}

// Reset clears all counts.
func (s *Sketch) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.table {
		for j := range s.table[i] {
			s.table[i][j] = 0
		}
	}
}

// TotalCount returns the sum of all deltas ever added to the sketch.
// Because each insertion updates every row identically, the total is
// read from row 0 by summing its buckets. This reflects the total
// event volume, not the number of distinct keys.
func (s *Sketch) TotalCount() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total uint64
	for _, v := range s.table[0] {
		total += v
	}
	return total
}

// bucket returns the column index for key in row i using a seeded FNV hash.
func (s *Sketch) bucket(key string, seed int) uint64 {
	h := fnv.New64a()
	// mix the seed into the hash to produce independent rows
	_, _ = h.Write([]byte{byte(seed), byte(seed >> 8)})
	_, _ = h.Write([]byte(key))
	return h.Sum64() % s.width
}
