package digest_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/digest"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Message:   "hello world",
		Level:     "info",
		Service:   "api",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Extra:     map[string]any{"request_id": "abc123", "status": 200},
	}
}

func TestNewEmptyOptionsReturnsError(t *testing.T) {
	_, err := digest.New(digest.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestSumIsConsistent(t *testing.T) {
	d, err := digest.New(digest.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := baseEntry()
	if got1, got2 := d.Sum(e), d.Sum(e); got1 != got2 {
		t.Fatalf("expected identical digests, got %q and %q", got1, got2)
	}
}

func TestSumDiffersOnMessageChange(t *testing.T) {
	d, _ := digest.New(digest.DefaultOptions())
	e1 := baseEntry()
	e2 := baseEntry()
	e2.Message = "different message"
	if d.Sum(e1) == d.Sum(e2) {
		t.Fatal("expected different digests for different messages")
	}
}

func TestSumDiffersOnLevelChange(t *testing.T) {
	d, _ := digest.New(digest.DefaultOptions())
	e1 := baseEntry()
	e2 := baseEntry()
	e2.Level = "error"
	if d.Sum(e1) == d.Sum(e2) {
		t.Fatal("expected different digests for different levels")
	}
}

func TestSumWithExtraKeys(t *testing.T) {
	d, err := digest.New(digest.Options{
		IncludeMessage: true,
		ExtraKeys:      []string{"request_id"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e1 := baseEntry()
	e2 := baseEntry()
	e2.Extra = map[string]any{"request_id": "xyz999", "status": 200}
	if d.Sum(e1) == d.Sum(e2) {
		t.Fatal("expected different digests for different extra values")
	}
}

func TestSumIgnoresTimestamp(t *testing.T) {
	d, _ := digest.New(digest.DefaultOptions())
	e1 := baseEntry()
	e2 := baseEntry()
	e2.Timestamp = time.Now().Add(24 * time.Hour)
	if d.Sum(e1) != d.Sum(e2) {
		t.Fatal("expected identical digests regardless of timestamp")
	}
}

func TestSumLengthIs64Chars(t *testing.T) {
	d, _ := digest.New(digest.DefaultOptions())
	s := d.Sum(baseEntry())
	if len(s) != 64 {
		t.Fatalf("expected 64-char hex digest, got %d chars", len(s))
	}
}

func TestSumMissingExtraKeyIgnored(t *testing.T) {
	d, _ := digest.New(digest.Options{
		IncludeMessage: true,
		ExtraKeys:      []string{"nonexistent"},
	})
	e := baseEntry()
	// Should not panic and should return a valid digest.
	s := d.Sum(e)
	if len(s) != 64 {
		t.Fatalf("expected 64-char hex digest, got %d", len(s))
	}
}
