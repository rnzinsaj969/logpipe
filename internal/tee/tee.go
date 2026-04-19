// Package tee provides a processor that forwards a log entry to multiple
// downstream handlers simultaneously, without modifying the entry.
package tee

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Sink is any function that accepts a log entry.
type Sink func(entry reader.LogEntry)

// Tee fans a single entry out to one or more sinks.
type Tee struct {
	sinks []Sink
}

// New returns a Tee that will forward entries to all provided sinks.
// At least one sink must be supplied.
func New(sinks ...Sink) (*Tee, error) {
	if len(sinks) == 0 {
		return nil, errors.New("tee: at least one sink is required")
	}
	for i, s := range sinks {
		if s == nil {
			return nil, errors.New("tee: sink at index " + itoa(i) + " is nil")
		}
	}
	return &Tee{sinks: sinks}, nil
}

// Apply sends entry to every registered sink in registration order.
func (t *Tee) Apply(entry reader.LogEntry) {
	for _, s := range t.sinks {
		s(entry)
	}
}

// Len returns the number of registered sinks.
func (t *Tee) Len() int { return len(t.sinks) }

// itoa is a minimal int-to-string helper to avoid importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
