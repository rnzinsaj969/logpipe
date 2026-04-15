// Package split provides a processor that splits a single log entry into
// multiple entries by expanding a repeated field into individual records.
package split

import (
	"errors"
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls the behaviour of the Splitter.
type Options struct {
	// Field is the key inside LogEntry.Extra whose value is a []any slice.
	// Each element of the slice becomes a separate LogEntry.
	Field string

	// KeepOriginal, when true, also emits the original entry unchanged.
	KeepOriginal bool
}

// Splitter expands one log entry into many based on a repeated extra field.
type Splitter struct {
	opts Options
}

// New returns a Splitter configured with opts.
// It returns an error if opts.Field is empty.
func New(opts Options) (*Splitter, error) {
	if opts.Field == "" {
		return nil, errors.New("split: Field must not be empty")
	}
	return &Splitter{opts: opts}, nil
}

// Apply expands entry into one or more entries.
//
// If the target field is absent or is not a []any, Apply returns the
// original entry wrapped in a single-element slice.
// Otherwise it returns one entry per element; each entry is a shallow copy
// of the original with the target field replaced by the element value.
// If KeepOriginal is true the original entry is prepended to the result.
func (s *Splitter) Apply(entry reader.LogEntry) []reader.LogEntry {
	raw, ok := entry.Extra[s.opts.Field]
	if !ok {
		return []reader.LogEntry{entry}
	}

	slice, ok := raw.([]any)
	if !ok {
		return []reader.LogEntry{entry}
	}

	var out []reader.LogEntry
	if s.opts.KeepOriginal {
		out = append(out, entry)
	}

	for i, elem := range slice {
		cloned := cloneEntry(entry)
		cloned.Extra[s.opts.Field] = elem
		cloned.Message = fmt.Sprintf("%s [%d]", entry.Message, i)
		out = append(out, cloned)
	}
	return out
}

// ApplyAll applies the splitter to each entry in entries and returns the
// concatenated results. It is a convenience wrapper around Apply for
// processing a batch of log entries in one call.
func (s *Splitter) ApplyAll(entries []reader.LogEntry) []reader.LogEntry {
	var out []reader.LogEntry
	for _, e := range entries {
		out = append(out, s.Apply(e)...)
	}
	return out
}

// cloneEntry returns a shallow copy of e with a new Extra map.
func cloneEntry(e reader.LogEntry) reader.LogEntry {
	cloned := e
	cloned.Extra = make(map[string]any, len(e.Extra))
	for k, v := range e.Extra {
		cloned.Extra[k] = v
	}
	return cloned
}
