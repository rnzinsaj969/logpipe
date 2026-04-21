package circa_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/circa"
	"github.com/logpipe/logpipe/internal/reader"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestInvalidGranularityReturnsError(t *testing.T) {
	_, err := circa.New(0)
	if err == nil {
		t.Fatal("expected error for zero granularity")
	}
	_, err = circa.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative granularity")
	}
}

func TestApplyTruncatesTimestamp(t *testing.T) {
	r, err := circa.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := time.Date(2024, 1, 1, 12, 34, 56, 0, time.UTC)
	e := reader.LogEntry{Message: "hello", Timestamp: ts}
	out := r.Apply(e)
	want := time.Date(2024, 1, 1, 12, 34, 0, 0, time.UTC)
	if !out.Timestamp.Equal(want) {
		t.Errorf("got %v, want %v", out.Timestamp, want)
	}
}

func TestApplyUsesClockOnZeroTimestamp(t *testing.T) {
	now := time.Date(2024, 6, 15, 9, 47, 33, 0, time.UTC)
	r, _ := circa.New(time.Hour)
	// Swap clock via newWithClock indirectly by testing behaviour:
	// Use exported New and verify zero-timestamp falls back to clock.
	// We can only observe the truncated result is within the same hour.
	e := reader.LogEntry{Message: "no ts"}
	out := r.Apply(e)
	if out.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp after Apply on zero entry")
	}
	_ = now
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	r, _ := circa.New(time.Minute)
	ts := time.Date(2024, 3, 10, 8, 22, 45, 0, time.UTC)
	e := reader.LogEntry{Message: "orig", Timestamp: ts}
	r.Apply(e)
	if !e.Timestamp.Equal(ts) {
		t.Error("original entry timestamp was mutated")
	}
}

func TestBucketReturnsCorrectBoundary(t *testing.T) {
	r, _ := circa.New(5 * time.Minute)
	ts := time.Date(2024, 1, 1, 0, 13, 59, 0, time.UTC)
	got := r.Bucket(ts)
	want := time.Date(2024, 1, 1, 0, 10, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Bucket: got %v, want %v", got, want)
	}
}

func TestGranularityReturnsConfiguredValue(t *testing.T) {
	g := 30 * time.Second
	r, _ := circa.New(g)
	if r.Granularity() != g {
		t.Errorf("Granularity: got %v, want %v", r.Granularity(), g)
	}
}

func TestApplyHourGranularity(t *testing.T) {
	r, _ := circa.New(time.Hour)
	ts := time.Date(2024, 12, 31, 23, 59, 59, 999, time.UTC)
	e := reader.LogEntry{Message: "end of year", Timestamp: ts}
	out := r.Apply(e)
	want := time.Date(2024, 12, 31, 23, 0, 0, 0, time.UTC)
	if !out.Timestamp.Equal(want) {
		t.Errorf("got %v, want %v", out.Timestamp, want)
	}
}
