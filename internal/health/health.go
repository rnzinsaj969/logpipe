package health

import (
	"sync"
	"time"
)

// Status represents the health state of a source.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
	StatusDown    Status = "down"
)

// SourceHealth holds health information for a single log source.
type SourceHealth struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	LastSeen  time.Time `json:"last_seen"`
	ErrorCount int      `json:"error_count"`
}

// Monitor tracks health state across multiple sources.
type Monitor struct {
	mu      sync.RWMutex
	sources map[string]*SourceHealth
}

// New creates a new health Monitor.
func New() *Monitor {
	return &Monitor{
		sources: make(map[string]*SourceHealth),
	}
}

// RecordSuccess marks a source as healthy and updates its last-seen timestamp.
func (m *Monitor) RecordSuccess(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreate(name)
	s.Status = StatusOK
	s.LastSeen = time.Now()
}

// RecordError increments the error count for a source and updates its status.
func (m *Monitor) RecordError(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.getOrCreate(name)
	s.ErrorCount++
	if s.ErrorCount >= 5 {
		s.Status = StatusDown
	} else {
		s.Status = StatusDegraded
	}
}

// Snapshot returns a copy of all current source health states.
func (m *Monitor) Snapshot() []SourceHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]SourceHealth, 0, len(m.sources))
	for _, s := range m.sources {
		out = append(out, *s)
	}
	return out
}

// getOrCreate returns an existing SourceHealth or creates a new one.
// Caller must hold m.mu.
func (m *Monitor) getOrCreate(name string) *SourceHealth {
	if s, ok := m.sources[name]; ok {
		return s
	}
	s := &SourceHealth{Name: name, Status: StatusOK}
	m.sources[name] = s
	return s
}
