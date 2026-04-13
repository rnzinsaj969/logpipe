// Package reader provides functionality for reading and parsing
// structured log entries from various input sources.
package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single structured log record.
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Reader reads LogEntry values from an io.Reader line by line.
type Reader struct {
	scanner *bufio.Scanner
	source  string
}

// New creates a new Reader that reads from r.
// source is a human-readable label (e.g. service name or file path).
func New(r io.Reader, source string) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
		source:  source,
	}
}

// Next reads the next log line and returns a parsed LogEntry.
// Returns io.EOF when there are no more lines.
func (r *Reader) Next() (*LogEntry, error) {
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return nil, fmt.Errorf("reader %q: scan error: %w", r.source, err)
		}
		return nil, io.EOF
	}

	line := r.scanner.Bytes()
	entry, err := parseLine(line)
	if err != nil {
		return nil, fmt.Errorf("reader %q: parse error: %w", r.source, err)
	}

	if entry.Service == "" {
		entry.Service = r.source
	}

	return entry, nil
}

// parseLine decodes a JSON-encoded log line into a LogEntry.
func parseLine(line []byte) (*LogEntry, error) {
	var entry LogEntry
	if err := json.Unmarshal(line, &entry); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if entry.Message == "" {
		return nil, fmt.Errorf("missing required field: message")
	}
	return &entry, nil
}
