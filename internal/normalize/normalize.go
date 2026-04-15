package normalize

import (
	"strings"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls normalization behaviour.
type Options struct {
	// LowercaseLevel converts the level field to lowercase.
	LowercaseLevel bool
	// TrimSpace trims leading and trailing whitespace from message and service.
	TrimSpace bool
	// DefaultLevel is used when the level field is empty.
	DefaultLevel string
	// DefaultService is used when the service field is empty.
	DefaultService string
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		LowercaseLevel: true,
		TrimSpace:      true,
		DefaultLevel:   "info",
		DefaultService: "unknown",
	}
}

// Normalizer applies field normalisation to log entries.
type Normalizer struct {
	opts Options
}

// New creates a Normalizer with the given options.
func New(opts Options) *Normalizer {
	return &Normalizer{opts: opts}
}

// Apply returns a new LogEntry with normalised fields.
// The original entry is never mutated.
func (n *Normalizer) Apply(e reader.LogEntry) reader.LogEntry {
	out := e

	if n.opts.TrimSpace {
		out.Message = strings.TrimSpace(out.Message)
		out.Service = strings.TrimSpace(out.Service)
	}

	if n.opts.LowercaseLevel {
		out.Level = strings.ToLower(out.Level)
	}

	if out.Level == "" && n.opts.DefaultLevel != "" {
		out.Level = n.opts.DefaultLevel
	}

	if out.Service == "" && n.opts.DefaultService != "" {
		out.Service = n.opts.DefaultService
	}

	return out
}
