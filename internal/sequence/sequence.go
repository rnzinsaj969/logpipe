package sequence

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Processor is any type that can transform a log entry.
type Processor interface {
	Apply(entry reader.LogEntry) (reader.LogEntry, error)
}

// Sequence runs a log entry through an ordered chain of Processors.
// Each processor receives the output of the previous one. If any
// processor returns an error the chain is aborted and the error is
// returned to the caller.
type Sequence struct {
	steps []Processor
}

// New creates a Sequence from the supplied processors. At least one
// processor must be provided.
func New(steps ...Processor) (*Sequence, error) {
	if len(steps) == 0 {
		return nil, fmt.Errorf("sequence: at least one processor is required")
	}
	return &Sequence{steps: steps}, nil
}

// Apply passes entry through each processor in order and returns the
// final transformed entry.
func (s *Sequence) Apply(entry reader.LogEntry) (reader.LogEntry, error) {
	var err error
	for i, p := range s.steps {
		entry, err = p.Apply(entry)
		if err != nil {
			return reader.LogEntry{}, fmt.Errorf("sequence: step %d: %w", i, err)
		}
	}
	return entry, nil
}

// Len returns the number of processors in the sequence.
func (s *Sequence) Len() int { return len(s.steps) }
