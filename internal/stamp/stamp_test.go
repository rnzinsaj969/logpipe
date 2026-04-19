package stamp_test

import (
	"testing"
	"time"

	"logpipe/internal/reader"
	"logpipe/internal/stamp"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
var fixed = time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNilClockReturnsError(t *testing.T) {
	_, err := stamp.New(stamp.Options{Clock: nil})
	if err == nil {
		t.Fatal("expected error for nil Clock")
	}
}

func TestSetsTimestampOnZeroEntry(t *testing.T) {
	s, _ := stamp.New(stamp.Options{Clock: fixedClock(fixed)})
	e := reader.LogEntry{Message: "hello"}
	out := s.Apply(e)
	if !out.Timestamp.Equal(fixed) {
		t.Fatalf("expected %v got %v", fixed, out.Timestamp)
	}
}

func TestDoesNotOverwriteByDefault(t *testing.T) {
	s, _ := stamp.New(stamp.Options{Clock: fixedClock(fixed)})
	e := reader.LogEntry{Message: "hello", Timestamp: epoch}
	out := s.Apply(e)
	if !out.Timestamp.Equal(epoch) {
		t.Fatalf("expected original timestamp %v got %v", epoch, out.Timestamp)
	}
}

func TestOverwriteExistingTimestamp(t *testing.T) {
	s, _ := stamp.New(stamp.Options{Clock: fixedClock(fixed), OverwriteExisting: true})
	e := reader.LogEntry{Message: "hello", Timestamp: epoch}
	out := s.Apply(e)
	if !out.Timestamp.Equal(fixed) {
		t.Fatalf("expected %v got %v", fixed, out.Timestamp)
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	s, _ := stamp.New(stamp.Options{Clock: fixedClock(fixed)})
	e := reader.LogEntry{Message: "hello"}
	s.Apply(e)
	if !e.Timestamp.IsZero() {
		t.Fatal("original entry was mutated")
	}
}

func TestDefaultOptionsNilClockPanicsGracefully(t *testing.T) {
	opts := stamp.DefaultOptions()
	if opts.Clock == nil {
		t.Fatal("DefaultOptions should provide a non-nil Clock")
	}
	s, err := stamp.New(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := reader.LogEntry{Message: "test"}
	out := s.Apply(e)
	if out.Timestamp.IsZero() {
		t.Fatal("expected timestamp to be set")
	}
}
