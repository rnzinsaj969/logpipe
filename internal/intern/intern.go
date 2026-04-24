// Package intern provides a string interning pool that deduplicates
// repeated field values (e.g. level, service) to reduce allocations
// across high-volume log entry pipelines.
package intern

import "sync"

// Pool is a concurrency-safe string intern pool.
type Pool struct {
	mu    sync.RWMutex
	table map[string]string
}

// New returns an initialised Pool ready for use.
func New() *Pool {
	return &Pool{table: make(map[string]string)}
}

// Intern returns the canonical copy of s stored in the pool.
// If s has not been seen before it is added and returned as-is.
func (p *Pool) Intern(s string) string {
	p.mu.RLock()
	if v, ok := p.table[s]; ok {
		p.mu.RUnlock()
		return v
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	// Double-check after acquiring write lock.
	if v, ok := p.table[s]; ok {
		return v
	}
	p.table[s] = s
	return s
}

// Len returns the number of unique strings currently held in the pool.
func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.table)
}

// Reset clears all interned strings from the pool.
func (p *Pool) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.table = make(map[string]string)
}

// InternSlice interns every element of ss in-place and returns ss.
func (p *Pool) InternSlice(ss []string) []string {
	for i, s := range ss {
		ss[i] = p.Intern(s)
	}
	return ss
}
