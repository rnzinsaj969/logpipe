// Package proxy forwards log entries to one or more named sinks.
package proxy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Sink is a destination that can receive log entries.
type Sink interface {
	Write(entry reader.LogEntry) error
}

// Proxy routes entries to registered sinks.
type Proxy struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// New returns an empty Proxy.
func New() *Proxy {
	return &Proxy{sinks: make(map[string]Sink)}
}

// Register adds a named sink. Returns an error if the name is already taken.
func (p *Proxy) Register(name string, s Sink) error {
	if name == "" {
		return errors.New("proxy: sink name must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.sinks[name]; ok {
		return fmt.Errorf("proxy: sink %q already registered", name)
	}
	p.sinks[name] = s
	return nil
}

// Remove deletes a sink by name.
func (p *Proxy) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.sinks, name)
}

// Forward sends the entry to every registered sink. All errors are collected
// and returned as a single combined error.
func (p *Proxy) Forward(entry reader.LogEntry) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var errs []error
	for name, s := range p.sinks {
		if err := s.Write(entry); err != nil {
			errs = append(errs, fmt.Errorf("proxy: sink %q: %w", name, err))
		}
	}
	return errors.Join(errs...)
}

// Len returns the number of registered sinks.
func (p *Proxy) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.sinks)
}
