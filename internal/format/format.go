// Package format provides utilities for formatting log entry fields
// into human-readable or structured string representations.
package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

// Options controls how a log entry is formatted.
type Options struct {
	// TimeFormat is the Go time layout used for the Timestamp field.
	// Defaults to time.RFC3339 when empty.
	TimeFormat string

	// UpperLevel converts the level string to upper-case when true.
	UpperLevel bool

	// ExtraDelimiter separates extra key=value pairs in text output.
	// Defaults to " " (single space) when empty.
	ExtraDelimiter string
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		TimeFormat:     time.RFC3339,
		UpperLevel:     false,
		ExtraDelimiter: " ",
	}
}

// Formatter formats log entries into strings.
type Formatter struct {
	opts Options
}

// New creates a Formatter with the given options.
func New(opts Options) *Formatter {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}
	if opts.ExtraDelimiter == "" {
		opts.ExtraDelimiter = " "
	}
	return &Formatter{opts: opts}
}

// FormatText returns a single-line text representation of the entry.
// Format: <timestamp> [<level>] <service>: <message> [key=value ...]
func (f *Formatter) FormatText(e reader.LogEntry) string {
	level := e.Level
	if f.opts.UpperLevel {
		level = strings.ToUpper(level)
	}

	ts := e.Timestamp.Format(f.opts.TimeFormat)
	base := fmt.Sprintf("%s [%s] %s: %s", ts, level, e.Service, e.Message)

	if len(e.Extra) == 0 {
		return base
	}

	parts := make([]string, 0, len(e.Extra))
	for k, v := range e.Extra {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return base + f.opts.ExtraDelimiter + strings.Join(parts, f.opts.ExtraDelimiter)
}

// FormatLevel returns the level string according to the formatter's options.
func (f *Formatter) FormatLevel(level string) string {
	if f.opts.UpperLevel {
		return strings.ToUpper(level)
	}
	return level
}

// FormatTimestamp formats a time.Time value using the configured layout.
func (f *Formatter) FormatTimestamp(t time.Time) string {
	return t.Format(f.opts.TimeFormat)
}
