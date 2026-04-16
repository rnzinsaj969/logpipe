// Package scope provides field scoping for log entries,
// allowing processors to operate on a subset of extra fields
// under a named namespace.
package scope

import (
	"errors"
	"fmt"

	"logpipe/internal/reader"
)

// Scoper extracts and re-embeds a named namespace within Extra.
type Scoper struct {
	namespace string
}

// New returns a Scoper for the given namespace key.
// Returns an error if namespace is empty.
func New(namespace string) (*Scoper, error) {
	if namespace == "" {
		return nil, errors.New("scope: namespace must not be empty")
	}
	return &Scoper{namespace: namespace}, nil
}

// Extract returns the nested map stored at the namespace key within entry.Extra.
// If the key is absent or not a map, an empty map is returned.
func (s *Scoper) Extract(entry reader.LogEntry) map[string]any {
	if entry.Extra == nil {
		return map[string]any{}
	}
	v, ok := entry.Extra[s.namespace]
	if !ok {
		return map[string]any{}
	}
	m, ok := v.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	out := make(map[string]any, len(m))
	for k, val := range m {
		out[k] = val
	}
	return out
}

// Embed returns a new LogEntry with the provided map stored under the namespace
// key in Extra, replacing any previous value. The original entry is not mutated.
func (s *Scoper) Embed(entry reader.LogEntry, fields map[string]any) reader.LogEntry {
	extra := make(map[string]any)
	for k, v := range entry.Extra {
		extra[k] = v
	}
	scoped := make(map[string]any, len(fields))
	for k, v := range fields {
		scoped[k] = v
	}
	extra[s.namespace] = scoped
	return reader.LogEntry{
		Service:   entry.Service,
		Level:     entry.Level,
		Message:   entry.Message,
		Timestamp: entry.Timestamp,
		Extra:     extra,
	}
}

// Key returns the namespace key used by this Scoper.
func (s *Scoper) Key() string {
	return fmt.Sprintf("%s", s.namespace)
}
