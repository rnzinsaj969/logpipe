// Package reader provides utilities for reading and parsing structured log
// entries from various sources.
//
// The primary entry point is New, which wraps an io.Reader and exposes a
// line-by-line Next method returning parsed LogEntry values.
//
// Each line is expected to be a JSON object containing at minimum a "message"
// field. Optional well-known fields include "level", "service", and "time".
// Lines that cannot be decoded as JSON are skipped with an error recorded on
// the entry so callers can decide how to handle malformed input.
//
// NewFileReader builds a Reader backed by a file path, returning an error if
// the file cannot be opened.
package reader
