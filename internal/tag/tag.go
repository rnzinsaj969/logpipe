// Package tag provides a processor that attaches a fixed set of key-value
// string tags to every log entry's Extra map.
package tag

import (
	"errors"
	"fmt"

	"logpipe/internal/reader"
)

// Tagger appends a static set of tags to each log entry.
type Tagger struct {
	tags map[string]string
}

// New creates a Tagger from the provided key-value pairs.
// tags must contain an even number of elements (key, value, key, value, …).
// Returns an error if any key is empty or the slice length is odd.
func New(tags ...string) (*Tagger, error) {
	if len(tags)%2 != 0 {
		return nil, errors.New("tag: tags must be key-value pairs (odd number of arguments)")
	}
	m := make(map[string]string, len(tags)/2)
	for i := 0; i < len(tags); i += 2 {
		k, v := tags[i], tags[i+1]
		if k == "" {
			return nil, fmt.Errorf("tag: key at position %d must not be empty", i)
		}
		m[k] = v
	}
	return &Tagger{tags: m}, nil
}

// Apply returns a copy of e with all configured tags merged into Extra.
// Existing Extra keys are not overwritten.
func (t *Tagger) Apply(e reader.LogEntry) reader.LogEntry {
	extra := make(map[string]any, len(e.Extra)+len(t.tags))
	for k, v := range t.tags {
		extra[k] = v
	}
	// Entry's own Extra values win over tag defaults.
	for k, v := range e.Extra {
		extra[k] = v
	}
	e.Extra = extra
	return e
}

// Len returns the number of configured tags.
func (t *Tagger) Len() int { return len(t.tags) }
