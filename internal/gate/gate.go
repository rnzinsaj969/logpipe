// Package gate provides a conditional pass-through filter that enables or
// disables log entry flow based on a runtime boolean flag.
package gate

import (
	"errors"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Gate conditionally forwards log entries depending on whether it is open.
type Gate struct {
	mu     sync.RWMutex
	open   bool
}

// New returns a Gate. If open is true entries are forwarded; otherwise they
// are dropped.
func New(open bool) (*Gate, error) {
	return &Gate{open: open}, nil
}

// Open sets the gate to the open state so that entries are forwarded.
func (g *Gate) Open() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.open = true
}

// Close sets the gate to the closed state so that entries are dropped.
func (g *Gate) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.open = false
}

// IsOpen reports whether the gate is currently open.
func (g *Gate) IsOpen() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.open
}

// Apply returns the entry unchanged when the gate is open. When the gate is
// closed it returns a sentinel error so callers can skip the entry.
func (g *Gate) Apply(e reader.LogEntry) (reader.LogEntry, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if !g.open {
		return reader.LogEntry{}, errors.New("gate: entry dropped (gate closed)")
	}
	return e, nil
}
