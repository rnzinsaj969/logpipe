// Package observe provides a lightweight observer that attaches a read-only
// tap to a log entry stream, invoking a callback for every entry that passes
// through without modifying or blocking the stream.
package observe

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Handler is called for every entry seen by the Observer.
type Handler func(entry reader.LogEntry)

// Observer wraps a Handler and forwards entries unchanged.
type Observer struct {
	handler Handler
}

// New returns an Observer that calls h for every entry passed to Apply.
// h must not be nil.
func New(h Handler) (*Observer, error) {
	if h == nil {
		return nil, errors.New("observe: handler must not be nil")
	}
	return &Observer{handler: h}, nil
}

// Apply calls the handler with entry and returns the entry unchanged.
// It never returns an error so it is safe to use in any pipeline stage.
func (o *Observer) Apply(entry reader.LogEntry) (reader.LogEntry, error) {
	o.handler(entry)
	return entry, nil
}
