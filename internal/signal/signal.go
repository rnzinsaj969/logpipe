// Package signal provides a broadcast notification mechanism that allows
// multiple subscribers to be notified when a named signal is fired.
package signal

import (
	"fmt"
	"sync"
)

// Handler is a function invoked when a signal is fired.
type Handler func(name string, payload map[string]any)

// Bus manages named signal subscriptions and dispatch.
type Bus struct {
	mu       sync.RWMutex
	subs     map[string][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{
		subs: make(map[string][]Handler),
	}
}

// Subscribe registers h to be called whenever name is fired.
// Returns an error if name is empty or h is nil.
func (b *Bus) Subscribe(name string, h Handler) error {
	if name == "" {
		return fmt.Errorf("signal: name must not be empty")
	}
	if h == nil {
		return fmt.Errorf("signal: handler must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[name] = append(b.subs[name], h)
	return nil
}

// Fire dispatches payload to all handlers subscribed to name.
// It is a no-op when no handlers are registered.
func (b *Bus) Fire(name string, payload map[string]any) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subs[name]))
	copy(handlers, b.subs[name])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(name, payload)
	}
}

// Reset removes all subscriptions for name.
func (b *Bus) Reset(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs, name)
}

// Len returns the number of handlers registered for name.
func (b *Bus) Len(name string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs[name])
}
