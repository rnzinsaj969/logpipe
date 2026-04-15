package sequence_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/sequence"
)

// upperProcessor uppercases the message field.
type upperProcessor struct{}

func (u upperProcessor) Apply(e reader.LogEntry) (reader.LogEntry, error) {
	e.Message = strings.ToUpper(e.Message)
	return e, nil
}

// prefixProcessor prepends a fixed string to the service field.
type prefixProcessor struct{ prefix string }

func (p prefixProcessor) Apply(e reader.LogEntry) (reader.LogEntry, error) {
	e.Service = p.prefix + e.Service
	return e, nil
}

// errProcessor always returns an error.
type errProcessor struct{}

func (errProcessor) Apply(_ reader.LogEntry) (reader.LogEntry, error) {
	return reader.LogEntry{}, errors.New("boom")
}

func base() reader.LogEntry {
	return reader.LogEntry{Service: "svc", Level: "info", Message: "hello"}
}

func TestSequenceAppliesAllSteps(t *testing.T) {
	seq, err := sequence.New(upperProcessor{}, prefixProcessor{prefix: "app-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := seq.Apply(base())
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if out.Message != "HELLO" {
		t.Errorf("message: got %q, want %q", out.Message, "HELLO")
	}
	if out.Service != "app-svc" {
		t.Errorf("service: got %q, want %q", out.Service, "app-svc")
	}
}

func TestSequenceAbortsOnError(t *testing.T) {
	seq, err := sequence.New(errProcessor{}, upperProcessor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = seq.Apply(base())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "step 0") {
		t.Errorf("error should mention step index, got: %v", err)
	}
}

func TestSequenceEmptyReturnsError(t *testing.T) {
	_, err := sequence.New()
	if err == nil {
		t.Fatal("expected error for empty sequence")
	}
}

func TestSequenceLen(t *testing.T) {
	seq, _ := sequence.New(upperProcessor{}, prefixProcessor{}, upperProcessor{})
	if seq.Len() != 3 {
		t.Errorf("Len: got %d, want 3", seq.Len())
	}
}

func TestSequenceSingleStep(t *testing.T) {
	seq, err := sequence.New(upperProcessor{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := seq.Apply(base())
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if out.Message != "HELLO" {
		t.Errorf("got %q, want %q", out.Message, "HELLO")
	}
}
