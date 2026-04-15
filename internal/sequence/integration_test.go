package sequence_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/sequence"
)

// timestampProcessor sets a fixed timestamp when the entry has none.
type timestampProcessor struct{ ts time.Time }

func (p timestampProcessor) Apply(e reader.LogEntry) (reader.LogEntry, error) {
	if e.Timestamp.IsZero() {
		e.Timestamp = p.ts
	}
	return e, nil
}

func TestSequencePipelineIntegration(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	seq, err := sequence.New(
		timestampProcessor{ts: now},
		upperProcessor{},
		prefixProcessor{prefix: "prod-"},
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entries := []reader.LogEntry{
		{Service: "auth", Level: "error", Message: "token expired"},
		{Service: "api", Level: "info", Message: "request ok", Timestamp: time.Now()},
	}

	for _, e := range entries {
		out, err := seq.Apply(e)
		if err != nil {
			t.Fatalf("Apply(%q): %v", e.Service, err)
		}
		if out.Timestamp.IsZero() {
			t.Errorf("service %q: timestamp should be set", e.Service)
		}
		if out.Service == e.Service {
			t.Errorf("service %q: expected prefix applied", e.Service)
		}
	}
}

func TestSequenceDoesNotMutateInput(t *testing.T) {
	seq, _ := sequence.New(upperProcessor{})

	original := reader.LogEntry{Service: "svc", Level: "info", Message: "hello"}
	_, err := seq.Apply(original)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if original.Message != "hello" {
		t.Errorf("input mutated: got %q", original.Message)
	}
}
