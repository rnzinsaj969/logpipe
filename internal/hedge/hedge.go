// Package hedge provides a speculative execution guard that issues a
// duplicate request after a configurable delay, returning whichever
// result arrives first. When applied to log entries it forwards the
// entry immediately and, if a confirmation is not received within the
// hedge window, emits the entry a second time so downstream sinks are
// not starved by a slow consumer.
package hedge

import (
	"errors"
	"sync"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Sink is any function that accepts a log entry.
type Sink func(reader.LogEntry) error

// Hedger wraps a primary sink and a fallback sink. If the primary sink
// does not acknowledge the entry within Window, the fallback sink
// receives the same entry.
type Hedger struct {
	mu       sync.Mutex
	primary  Sink
	fallback Sink
	window   time.Duration
	clock    func() time.Time
	sleep    func(time.Duration)
}

// New returns a Hedger that forwards entries to primary and, after
// window elapses without a successful primary response, also forwards
// to fallback.
func New(primary, fallback Sink, window time.Duration) (*Hedger, error) {
	if primary == nil {
		return nil, errors.New("hedge: primary sink must not be nil")
	}
	if fallback == nil {
		return nil, errors.New("hedge: fallback sink must not be nil")
	}
	if window <= 0 {
		return nil, errors.New("hedge: window must be positive")
	}
	return &Hedger{
		primary:  primary,
		fallback: fallback,
		window:   window,
		clock:    time.Now,
		sleep:    time.Sleep,
	}, nil
}

// Apply sends entry to the primary sink. If the primary sink returns
// an error, or if it does not return within the hedge window, the
// fallback sink receives the same entry. The first non-nil error
// encountered is returned; a nil error is returned when at least one
// sink succeeds.
func (h *Hedger) Apply(entry reader.LogEntry) error {
	type result struct{ err error }
	ch := make(chan result, 1)

	go func() {
		ch <- result{err: h.primary(entry)}
	}()

	timer := time.NewTimer(h.window)
	defer timer.Stop()

	select {
	case res := <-ch:
		if res.err == nil {
			return nil
		}
		// Primary failed – try fallback synchronously.
		return h.fallback(entry)
	case <-timer.C:
		// Primary too slow – hedge with fallback.
		if err := h.fallback(entry); err != nil {
			// Wait for primary to finish before returning.
			res := <-ch
			if res.err != nil {
				return res.err
			}
			return nil
		}
		return nil
	}
}
