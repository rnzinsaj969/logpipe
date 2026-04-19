// Package fork splits a log entry stream into two independent branches
// based on a predicate function.
package fork

import (
	"errors"

	"github.com/logpipe/logpipe/internal/reader"
)

// Predicate returns true if the entry should go to the left branch.
type Predicate func(entry reader.LogEntry) bool

// Fork holds two sinks and a predicate.
type Fork struct {
	pred  Predicate
	left  func(reader.LogEntry) error
	right func(reader.LogEntry) error
}

// New creates a Fork. Entries matching pred are sent to left; others to right.
func New(pred Predicate, left, right func(reader.LogEntry) error) (*Fork, error) {
	if pred == nil {
		return nil, errors.New("fork: predicate must not be nil")
	}
	if left == nil {
		return nil, errors.New("fork: left sink must not be nil")
	}
	if right == nil {
		return nil, errors.New("fork: right sink must not be nil")
	}
	return &Fork{pred: pred, left: left, right: right}, nil
}

// Apply routes the entry to the appropriate sink.
func (f *Fork) Apply(entry reader.LogEntry) error {
	if f.pred(entry) {
		return f.left(entry)
	}
	return f.right(entry)
}
