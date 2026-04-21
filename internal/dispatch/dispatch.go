// Package dispatch routes log entries to named sinks based on a registry.
// Each sink is identified by a string key and receives entries via its Apply
// method. Dispatch is safe for concurrent use.
package dispatch

import (
	"errors"
	"fmt"
	"sync"

	"github.com/logpipe/logpipe/internal/reader"
)

// Sink is any value that can receive a log entry.
type Sink interface {
	Apply(entry reader.LogEntry) (reader.LogEntry, error)
}

// Dispatcher holds a registry of named sinks and routes entries to one or
// more of them by name.
type Dispatcher struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// New returns an empty Dispatcher.
func New() *Dispatcher {
	return &Dispatcher{
		sinks: make(map[string]Sink),
	}
}

// Register adds a sink under the given name. It returns an error if the name
// is empty, the sink is nil, or the name is already registered.
func (d *Dispatcher) Register(name string, sink Sink) error {
	if name == "" {
		return errors.New("dispatch: sink name must not be empty")
	}
	if sink == nil {
		return errors.New("dispatch: sink must not be nil")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.sinks[name]; exists {
		return fmt.Errorf("dispatch: sink %q already registered", name)
	}
	d.sinks[name] = sink
	return nil
}

// Unregister removes the sink with the given name. It is a no-op if the name
// is not registered.
func (d *Dispatcher) Unregister(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.sinks, name)
}

// Send delivers entry to the sink identified by name. It returns an error if
// the name is not registered or the sink returns an error.
func (d *Dispatcher) Send(name string, entry reader.LogEntry) (reader.LogEntry, error) {
	d.mu.RLock()
	sink, ok := d.sinks[name]
	d.mu.RUnlock()
	if !ok {
		return entry, fmt.Errorf("dispatch: no sink registered for %q", name)
	}
	return sink.Apply(entry)
}

// Broadcast delivers entry to every registered sink. Errors from individual
// sinks are collected and returned as a combined error; all sinks are called
// regardless of earlier failures.
func (d *Dispatcher) Broadcast(entry reader.LogEntry) error {
	d.mu.RLock()
	names := make([]string, 0, len(d.sinks))
	for name := range d.sinks {
		names = append(names, name)
	}
	d.mu.RUnlock()

	var errs []error
	for _, name := range names {
		d.mu.RLock()
		sink, ok := d.sinks[name]
		d.mu.RUnlock()
		if !ok {
			continue
		}
		if _, err := sink.Apply(entry); err != nil {
			errs = append(errs, fmt.Errorf("dispatch: sink %q: %w", name, err))
		}
	}
	return errors.Join(errs...)
}

// Names returns the registered sink names in an unspecified order.
func (d *Dispatcher) Names() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]string, 0, len(d.sinks))
	for name := range d.sinks {
		out = append(out, name)
	}
	return out
}
