// Package output handles formatting and writing log entries to various destinations.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Format represents the output format for log entries.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Entry represents a structured log entry for output.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// Writer writes log entries to an io.Writer in a specified format.
type Writer struct {
	out    io.Writer
	format Format
}

// New creates a new Writer with the given destination and format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Writer {
	if out == nil {
		out = os.Stdout
	}
	return &Writer{out: out, format: format}
}

// Write formats and writes a single log entry to the underlying writer.
func (w *Writer) Write(e Entry) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(e)
	default:
		return w.writeText(e)
	}
}

func (w *Writer) writeText(e Entry) error {
	ts := e.Timestamp.Format(time.RFC3339)
	_, err := fmt.Fprintf(w.out, "%s [%s] (%s) %s\n", ts, e.Level, e.Service, e.Message)
	return err
}

func (w *Writer) writeJSON(e Entry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("output: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(w.out, "%s\n", data)
	return err
}
