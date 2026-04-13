package reader

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestReaderNext(t *testing.T) {
	input := `{"timestamp":"2024-01-15T10:00:00Z","level":"info","service":"api","message":"server started"}
{"timestamp":"2024-01-15T10:00:01Z","level":"error","service":"api","message":"connection failed"}
`
	r := New(strings.NewReader(input), "api")

	entry1, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error on first entry: %v", err)
	}
	if entry1.Message != "server started" {
		t.Errorf("expected message %q, got %q", "server started", entry1.Message)
	}
	if entry1.Level != "info" {
		t.Errorf("expected level %q, got %q", "info", entry1.Level)
	}

	entry2, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error on second entry: %v", err)
	}
	if entry2.Level != "error" {
		t.Errorf("expected level %q, got %q", "error", entry2.Level)
	}

	_, err = r.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReaderFallbackService(t *testing.T) {
	input := `{"level":"warn","message":"disk usage high"}
`
	r := New(strings.NewReader(input), "monitor")

	entry, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Service != "monitor" {
		t.Errorf("expected service %q, got %q", "monitor", entry.Service)
	}
}

func TestReaderInvalidJSON(t *testing.T) {
	input := "not valid json\n"
	r := New(strings.NewReader(input), "svc")

	_, err := r.Next()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestReaderMissingMessage(t *testing.T) {
	input := `{"level":"info","service":"svc"}` + "\n"
	r := New(strings.NewReader(input), "svc")

	_, err := r.Next()
	if err == nil {
		t.Fatal("expected error for missing message field, got nil")
	}
}

func TestLogEntryTimestamp(t *testing.T) {
	input := `{"timestamp":"2024-06-01T12:00:00Z","level":"debug","message":"ping"}` + "\n"
	r := New(strings.NewReader(input), "svc")

	entry, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entry.Timestamp)
	}
}
