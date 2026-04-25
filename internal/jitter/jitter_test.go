package jitter_test

import (
	"testing"
	"time"

	"logpipe/internal/jitter"
	"logpipe/internal/reader"
)

// fixedSource always returns the same value so tests are deterministic.
type fixedSource struct{ val int64 }

func (f fixedSource) Int63n(_ int64) int64 { return f.val }

func base(ts time.Time) reader.LogEntry {
	return reader.LogEntry{Service: "svc", Level: "info", Message: "hello", Timestamp: ts}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := jitter.New(0, fixedSource{})
	if err == nil {
		t.Fatal("expected error for zero max")
	}
	_, err = jitter.New(-time.Second, fixedSource{})
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestApplyAddsOffset(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	j, err := jitter.New(time.Second, fixedSource{val: 250_000_000}) // 250ms
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := j.Apply(base(now))
	want := now.Add(250 * time.Millisecond)
	if !out.Timestamp.Equal(want) {
		t.Errorf("got %v, want %v", out.Timestamp, want)
	}
}

func TestApplyZeroTimestampUnchanged(t *testing.T) {
	j, _ := jitter.New(time.Second, fixedSource{val: 500_000_000})
	e := base(time.Time{})
	out := j.Apply(e)
	if !out.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp, got %v", out.Timestamp)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	now := time.Now()
	j, _ := jitter.New(time.Second, fixedSource{val: 1})
	e := base(now)
	_ = j.Apply(e)
	if !e.Timestamp.Equal(now) {
		t.Error("original entry was mutated")
	}
}

func TestApplyPreservesOtherFields(t *testing.T) {
	now := time.Now()
	j, _ := jitter.New(time.Millisecond, fixedSource{val: 0})
	e := base(now)
	e.Extra = map[string]any{"k": "v"}
	out := j.Apply(e)
	if out.Service != e.Service || out.Level != e.Level || out.Message != e.Message {
		t.Error("non-timestamp fields were altered")
	}
	if out.Extra["k"] != "v" {
		t.Error("extra fields were altered")
	}
}

func TestNilSourceUsesDefault(t *testing.T) {
	j, err := jitter.New(time.Millisecond, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Now()
	out := j.Apply(base(now))
	// Timestamp should be >= now and < now+1ms.
	if out.Timestamp.Before(now) || !out.Timestamp.Before(now.Add(time.Millisecond)) {
		t.Errorf("timestamp %v out of expected range", out.Timestamp)
	}
}
