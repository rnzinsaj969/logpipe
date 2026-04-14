// Package enrich provides log entry enrichment by attaching static or
// dynamic metadata fields to each entry as it passes through the pipeline.
//
// Enrichers are composable: multiple enrichers can be chained together,
// each adding its own set of fields without overwriting prior work unless
// explicitly configured to do so.
package enrich

import (
	"strings"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Field represents a single metadata key/value pair to attach to a log entry.
type Field struct {
	Key   string
	Value string
}

// Source is a function that returns a value for a dynamic field at the time
// an entry is processed. It receives the current entry for context.
type Source func(entry reader.LogEntry) string

// rule is an internal enrichment rule that may be static or dynamic.
type rule struct {
	key       string
	staticVal string
	dynamic   Source
	overwrite bool
}

// Enricher attaches metadata fields to log entries.
type Enricher struct {
	rules []rule
}

// New creates an Enricher with no rules. Use Add and AddDynamic to register
// enrichment rules before calling Apply.
func New() *Enricher {
	return &Enricher{}
}

// Add registers a static field. If overwrite is false the field is skipped
// when the entry already contains a non-empty value for the key.
func (e *Enricher) Add(key, value string, overwrite bool) *Enricher {
	e.rules = append(e.rules, rule{
		key:       key,
		staticVal: value,
		overwrite: overwrite,
	})
	return e
}

// AddDynamic registers a field whose value is computed from the entry at
// processing time. If overwrite is false the field is skipped when a
// non-empty value already exists.
func (e *Enricher) AddDynamic(key string, src Source, overwrite bool) *Enricher {
	e.rules = append(e.rules, rule{
		key:       key,
		dynamic:   src,
		overwrite: overwrite,
	})
	return e
}

// HasRules reports whether any enrichment rules have been registered.
func (e *Enricher) HasRules() bool {
	return len(e.rules) > 0
}

// Apply returns a copy of entry with all registered enrichment rules applied.
// The original entry is never mutated.
func (e *Enricher) Apply(entry reader.LogEntry) reader.LogEntry {
	out := entry
	if out.Extra == nil {
		out.Extra = make(map[string]string)
	} else {
		// shallow-copy so we do not mutate the caller's map
		copy := make(map[string]string, len(entry.Extra))
		for k, v := range entry.Extra {
			copy[k] = v
		}
		out.Extra = copy
	}

	for _, r := range e.rules {
		val := r.staticVal
		if r.dynamic != nil {
			val = r.dynamic(entry)
		}

		if !r.overwrite {
			if existing, ok := out.Extra[r.key]; ok && strings.TrimSpace(existing) != "" {
				continue
			}
		}

		out.Extra[r.key] = val
	}

	return out
}

// WithHostname returns a Source that always returns the provided hostname.
// Useful for stamping entries with the originating host.
func WithHostname(hostname string) Source {
	return func(_ reader.LogEntry) string { return hostname }
}

// WithTimestampUTC returns a Source that formats the entry's timestamp as an
// RFC3339 UTC string, falling back to the current time when the entry has a
// zero timestamp.
func WithTimestampUTC() Source {
	return func(entry reader.LogEntry) string {
		t := entry.Timestamp
		if t.IsZero() {
			t = time.Now().UTC()
		}
		return t.UTC().Format(time.RFC3339)
	}
}
