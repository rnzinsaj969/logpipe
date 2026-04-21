// Package circa provides approximate time-bucketing for log entries,
// rounding timestamps to a configurable granularity for grouping or
// comparison purposes.
package circa

import (
	"fmt"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

// Rounder rounds a log entry's timestamp to the nearest bucket boundary.
type Rounder struct {
	granularity time.Duration
	clock       func() time.Time
}

// New returns a Rounder that truncates entry timestamps to the given
// granularity. granularity must be greater than zero.
func New(granularity time.Duration) (*Rounder, error) {
	return newWithClock(granularity, time.Now)
}

func newWithClock(granularity time.Duration, clock func() time.Time) (*Rounder, error) {
	if granularity <= 0 {
		return nil, fmt.Errorf("circa: granularity must be greater than zero, got %s", granularity)
	}
	return &Rounder{granularity: granularity, clock: clock}, nil
}

// Apply returns a copy of e with its timestamp truncated to the Rounder's
// granularity. If e.Timestamp is the zero value the current clock time is
// used before truncation.
func (r *Rounder) Apply(e reader.LogEntry) reader.LogEntry {
	ts := e.Timestamp
	if ts.IsZero() {
		ts = r.clock()
	}
	out := e
	out.Timestamp = ts.Truncate(r.granularity)
	return out
}

// Bucket returns the bucket boundary time for an arbitrary timestamp.
func (r *Rounder) Bucket(t time.Time) time.Time {
	return t.Truncate(r.granularity)
}

// Granularity returns the configured bucket size.
func (r *Rounder) Granularity() time.Duration {
	return r.granularity
}
