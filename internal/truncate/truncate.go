// Package truncate provides utilities for capping log entry field lengths
// to prevent oversized payloads from propagating through the pipeline.
package truncate

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls which fields are truncated and at what byte limits.
type Options struct {
	MaxMessageBytes int
	MaxFieldBytes   int
}

// DefaultOptions returns sensible production defaults.
func DefaultOptions() Options {
	return Options{
		MaxMessageBytes: 4096,
		MaxFieldBytes:   512,
	}
}

// Truncator applies length caps to log entries.
type Truncator struct {
	opts Options
}

// New creates a Truncator with the given options.
func New(opts Options) *Truncator {
	if opts.MaxMessageBytes <= 0 {
		opts.MaxMessageBytes = DefaultOptions().MaxMessageBytes
	}
	if opts.MaxFieldBytes <= 0 {
		opts.MaxFieldBytes = DefaultOptions().MaxFieldBytes
	}
	return &Truncator{opts: opts}
}

// Apply returns a copy of the entry with oversized fields truncated.
// Truncated values are suffixed with "…" to indicate data loss.
func (t *Truncator) Apply(e reader.LogEntry) reader.LogEntry {
	out := reader.LogEntry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Service:   e.Service,
		Message:   truncateString(e.Message, t.opts.MaxMessageBytes),
		Fields:    make(map[string]any, len(e.Fields)),
	}
	for k, v := range e.Fields {
		if s, ok := v.(string); ok {
			out.Fields[k] = truncateString(s, t.opts.MaxFieldBytes)
		} else {
			out.Fields[k] = v
		}
	}
	return out
}

// Needed reports whether the entry contains any field that exceeds the limits.
func (t *Truncator) Needed(e reader.LogEntry) bool {
	if len(e.Message) > t.opts.MaxMessageBytes {
		return true
	}
	for _, v := range e.Fields {
		if s, ok := v.(string); ok && len(s) > t.opts.MaxFieldBytes {
			return true
		}
	}
	return false
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	// Cut on a rune boundary then append ellipsis.
	cut := []rune(s)
	if len(cut) > max {
		cut = cut[:max]
	}
	return fmt.Sprintf("%s…", string(cut))
}
