// Package flatten provides utilities for flattening nested log entry
// extra fields into a single-level map using dot-separated keys.
package flatten

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls flattening behaviour.
type Options struct {
	// Separator is placed between parent and child keys. Defaults to ".".
	Separator string
	// MaxDepth limits recursion. Zero means unlimited.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: ".",
		MaxDepth:  0,
	}
}

// Flattener flattens nested maps inside LogEntry.Extra.
type Flattener struct {
	opts Options
}

// New creates a Flattener with the provided options.
func New(opts Options) *Flattener {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	return &Flattener{opts: opts}
}

// Apply returns a new LogEntry whose Extra map has been flattened.
// The original entry is not mutated.
func (f *Flattener) Apply(e reader.LogEntry) reader.LogEntry {
	if len(e.Extra) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Extra))
	f.flattenMap("", e.Extra, out, 1)
	e.Extra = out
	return e
}

func (f *Flattener) flattenMap(prefix string, src map[string]any, dst map[string]any, depth int) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = prefix + f.opts.Separator + k
		}
		if nested, ok := v.(map[string]any); ok && (f.opts.MaxDepth == 0 || depth < f.opts.MaxDepth) {
			f.flattenMap(key, nested, dst, depth+1)
		} else {
			dst[key] = v
		}
	}
}

// FlatKey builds a dot-separated key from parts, useful for constructing
// expected keys in tests or downstream consumers.
func FlatKey(sep string, parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out = fmt.Sprintf("%s%s%s", out, sep, p)
	}
	return out
}
