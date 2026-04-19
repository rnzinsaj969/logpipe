// Package fanout distributes a single log entry to multiple independent sinks.
package fanout

import (
	"errors"
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Sink is any function that accepts a log entry.
type Sink func(reader.LogEntry) error

// Fanout sends each entry to all registered sinks.
type Fanout struct {
	sinks []Sink
}

// New creates a Fanout with the provided sinks.
// Returns an error if no sinks are provided or any sink is nil.
func New(sinks ...Sink) (*Fanout, error) {
	if len(sinks) == 0 {
		return nil, errors.New("fanout: at least one sink is required")
	}
	for i, s := range sinks {
		if s == nil {
			return nil, fmt.Errorf("fanout: sink at index %d is nil", i)
		}
	}
	return &Fanout{sinks: sinks}, nil
}

// Apply forwards entry to every sink. All sinks are called even if one errors.
// Returns a combined error if any sink fails.
func (f *Fanout) Apply(entry reader.LogEntry) error {
	var errs []error
	for _, s := range f.sinks {
		if err := s(entry); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// Len returns the number of registered sinks.
func (f *Fanout) Len() int { return len(f.sinks) }
