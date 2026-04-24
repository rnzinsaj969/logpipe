// Package admit provides a probabilistic admission filter that accepts
// log entries based on a configurable acceptance rate, using a fast
// hash of the entry message to ensure deterministic behaviour for
// identical messages within the same run.
package admit

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/logpipe/logpipe/internal/reader"
)

// Admitter accepts or rejects log entries based on a rate in [0, 1].
type Admitter struct {
	threshold uint32
}

// New returns an Admitter that admits entries with probability rate.
// rate is clamped to [0, 1]. A rate of 1 admits everything; 0 admits nothing.
func New(rate float64) (*Admitter, error) {
	if math.IsNaN(rate) {
		return nil, fmt.Errorf("admit: rate must not be NaN")
	}
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	threshold := uint32(rate * float64(math.MaxUint32))
	return &Admitter{threshold: threshold}, nil
}

// Admit returns true if the entry should be accepted.
// The decision is deterministic for a given message string.
func (a *Admitter) Admit(e reader.LogEntry) bool {
	if a.threshold == math.MaxUint32 {
		return true
	}
	if a.threshold == 0 {
		return false
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(e.Message))
	return h.Sum32() <= a.threshold
}

// Apply returns a new slice containing only admitted entries.
func (a *Admitter) Apply(entries []reader.LogEntry) []reader.LogEntry {
	out := make([]reader.LogEntry, 0, len(entries))
	for _, e := range entries {
		if a.Admit(e) {
			out = append(out, e)
		}
	}
	return out
}
