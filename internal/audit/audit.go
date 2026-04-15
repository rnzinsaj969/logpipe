// Package audit provides a structured audit trail for log pipeline events,
// recording when entries are dropped, transformed, or routed.
package audit

import (
	"sync"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventDropped     EventKind = "dropped"
	EventTransformed EventKind = "transformed"
	EventRouted      EventKind = "routed"
	EventRedacted    EventKind = "redacted"
)

// Event represents a single audit record.
type Event struct {
	Kind      EventKind         `json:"kind"`
	Service   string            `json:"service"`
	Reason    string            `json:"reason"`
	Timestamp time.Time         `json:"timestamp"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Log is an in-memory audit log with a configurable capacity.
type Log struct {
	mu       sync.Mutex
	events   []Event
	capacity int
}

// New creates a new audit Log with the given maximum capacity.
// When the capacity is reached the oldest event is evicted.
func New(capacity int) *Log {
	if capacity <= 0 {
		capacity = 256
	}
	return &Log{
		events:   make([]Event, 0, capacity),
		capacity: capacity,
	}
}

// Record appends an event to the audit log.
func (l *Log) Record(kind EventKind, service, reason string, meta map[string]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.events) >= l.capacity {
		l.events = l.events[1:]
	}
	l.events = append(l.events, Event{
		Kind:      kind,
		Service:   service,
		Reason:    reason,
		Timestamp: time.Now().UTC(),
		Meta:      meta,
	})
}

// Snapshot returns a copy of all current audit events.
func (l *Log) Snapshot() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Len returns the number of events currently stored.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.events)
}

// Clear removes all events from the log.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = l.events[:0]
}
