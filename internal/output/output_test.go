package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/output"
)

var testEntry = output.Entry{
	Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	Service:   "api",
	Level:     "INFO",
	Message:   "request received",
}

func TestWriteText(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatText)

	if err := w.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[INFO]") {
		t.Errorf("expected level in output, got: %s", got)
	}
	if !strings.Contains(got, "(api)") {
		t.Errorf("expected service in output, got: %s", got)
	}
	if !strings.Contains(got, "request received") {
		t.Errorf("expected message in output, got: %s", got)
	}
	if !strings.Contains(got, "2024-01-15T10:30:00Z") {
		t.Errorf("expected timestamp in output, got: %s", got)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON)

	if err := w.Write(testEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result output.Entry
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if result.Service != testEntry.Service {
		t.Errorf("service: got %q, want %q", result.Service, testEntry.Service)
	}
	if result.Level != testEntry.Level {
		t.Errorf("level: got %q, want %q", result.Level, testEntry.Level)
	}
	if result.Message != testEntry.Message {
		t.Errorf("message: got %q, want %q", result.Message, testEntry.Message)
	}
}

func TestWriteDefaultsToStdout(t *testing.T) {
	// Ensure New with nil out does not panic.
	w := output.New(nil, output.FormatText)
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}
