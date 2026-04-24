package admit

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(msg string) reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestAdmitRateOneAcceptsAll(t *testing.T) {
	a, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, msg := range []string{"alpha", "beta", "gamma", "delta"} {
		if !a.Admit(entry(msg)) {
			t.Errorf("expected entry %q to be admitted", msg)
		}
	}
}

func TestAdmitRateZeroRejectsAll(t *testing.T) {
	a, err := New(0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, msg := range []string{"alpha", "beta", "gamma", "delta"} {
		if a.Admit(entry(msg)) {
			t.Errorf("expected entry %q to be rejected", msg)
		}
	}
}

func TestAdmitRateClampedAboveOne(t *testing.T) {
	a, err := New(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.Admit(entry("any")) {
		t.Error("expected entry to be admitted after clamping rate to 1")
	}
}

func TestAdmitRateClampedBelowZero(t *testing.T) {
	a, err := New(-0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Admit(entry("any")) {
		t.Error("expected entry to be rejected after clamping rate to 0")
	}
}

func TestAdmitNaNReturnsError(t *testing.T) {
	_, err := New(math.NaN())
	if err == nil {
		t.Fatal("expected error for NaN rate")
	}
}

func TestAdmitDeterministic(t *testing.T) {
	a, err := New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := entry("stable-message")
	first := a.Admit(e)
	for i := 0; i < 10; i++ {
		if a.Admit(e) != first {
			t.Error("expected deterministic result for same message")
		}
	}
}

func TestApplyFiltersSlice(t *testing.T) {
	a, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := []reader.LogEntry{entry("a"), entry("b"), entry("c")}
	out := a.Apply(input)
	if len(out) != len(input) {
		t.Errorf("expected %d entries, got %d", len(input), len(out))
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	a, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := []reader.LogEntry{entry("x")}
	_ = a.Apply(input)
	if input[0].Message != "x" {
		t.Error("Apply must not mutate input entries")
	}
}
