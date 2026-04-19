package sieve_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/sieve"
)

func TestSieveFiltersSlice(t *testing.T) {
	s, err := sieve.New([]sieve.Rule{
		{Level: "error"},
		{Level: "warn", Service: "payments"},
	})
	if err != nil {
		t.Fatal(err)
	}

	input := []reader.LogEntry{
		entry("info", "api", "started"),
		entry("error", "api", "crashed"),
		entry("warn", "payments", "retry"),
		entry("warn", "api", "slow"),
		entry("debug", "worker", "tick"),
	}

	var kept []reader.LogEntry
	for _, e := range input {
		if s.Apply(e) {
			kept = append(kept, e)
		}
	}

	if len(kept) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(kept))
	}
	if kept[0].Message != "crashed" {
		t.Errorf("unexpected first entry: %s", kept[0].Message)
	}
	if kept[1].Message != "retry" {
		t.Errorf("unexpected second entry: %s", kept[1].Message)
	}
}
