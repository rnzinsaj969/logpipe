package expire_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/expire"
	"github.com/logpipe/logpipe/internal/reader"
)

var now = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return now }

func TestNewInvalidMaxAgeReturnsError(t *testing.T) {
	_, err := expire.New(expire.Options{MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
}

func TestApplyRetainsFreshEntry(t *testing.T) {
	p, _ := expire.New(expire.Options{MaxAge: time.Minute, Clock: fixedClock})
	e := reader.LogEntry{Timestamp: now.Add(-30 * time.Second)}
	if !p.Apply(e) {
		t.Error("expected fresh entry to be retained")
	}
}

func TestApplyDropsExpiredEntry(t *testing.T) {
	p, _ := expire.New(expire.Options{MaxAge: time.Minute, Clock: fixedClock})
	e := reader.LogEntry{Timestamp: now.Add(-2 * time.Minute)}
	if p.Apply(e) {
		t.Error("expected expired entry to be dropped")
	}
}

func TestApplyRetainsZeroTimestamp(t *testing.T) {
	p, _ := expire.New(expire.Options{MaxAge: time.Second, Clock: fixedClock})
	e := reader.LogEntry{} // zero timestamp
	if !p.Apply(e) {
		t.Error("expected zero-timestamp entry to pass through")
	}
}

func TestFilterRemovesExpiredEntries(t *testing.T) {
	p, _ := expire.New(expire.Options{MaxAge: time.Minute, Clock: fixedClock})
	entries := []reader.LogEntry{
		{Message: "fresh", Timestamp: now.Add(-10 * time.Second)},
		{Message: "old", Timestamp: now.Add(-5 * time.Minute)},
		{Message: "also fresh", Timestamp: now.Add(-59 * time.Second)},
	}
	got := p.Filter(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Message != "fresh" || got[1].Message != "also fresh" {
		t.Errorf("unexpected entries: %+v", got)
	}
}

func TestFilterAllExpired(t *testing.T) {
	p, _ := expire.New(expire.Options{MaxAge: time.Second, Clock: fixedClock})
	entries := []reader.LogEntry{
		{Timestamp: now.Add(-1 * time.Hour)},
		{Timestamp: now.Add(-2 * time.Hour)},
	}
	if got := p.Filter(entries); len(got) != 0 {
		t.Errorf("expected empty result, got %d entries", len(got))
	}
}
