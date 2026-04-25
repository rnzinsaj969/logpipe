package watermark_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/watermark"
)

func entry(ts time.Time) *reader.LogEntry {
	return &reader.LogEntry{Message: "msg", Level: "info", Service: "svc", Timestamp: ts}
}

func TestHighStartsAtZero(t *testing.T) {
	w := watermark.New()
	if !w.High().IsZero() {
		t.Fatalf("expected zero time, got %v", w.High())
	}
}

func TestAdvanceUpdatesHigh(t *testing.T) {
	w := watermark.New()
	now := time.Now()
	if err := w.Advance(entry(now)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.High().Equal(now) {
		t.Fatalf("expected %v, got %v", now, w.High())
	}
}

func TestAdvanceDoesNotGoBackwards(t *testing.T) {
	w := watermark.New()
	now := time.Now()
	_ = w.Advance(entry(now))
	_ = w.Advance(entry(now.Add(-time.Second)))
	if !w.High().Equal(now) {
		t.Fatalf("high-water mark moved backwards: got %v", w.High())
	}
}

func TestAdvanceNilEntryReturnsError(t *testing.T) {
	w := watermark.New()
	if err := w.Advance(nil); err == nil {
		t.Fatal("expected error for nil entry")
	}
}

func TestBehindReturnsTrueForLateEntry(t *testing.T) {
	w := watermark.New()
	now := time.Now()
	_ = w.Advance(entry(now))
	if !w.Behind(entry(now.Add(-time.Millisecond))) {
		t.Fatal("expected entry to be behind watermark")
	}
}

func TestBehindReturnsFalseForFreshEntry(t *testing.T) {
	w := watermark.New()
	now := time.Now()
	_ = w.Advance(entry(now))
	if w.Behind(entry(now.Add(time.Second))) {
		t.Fatal("expected entry to be ahead of watermark")
	}
}

func TestBehindNilReturnsFalse(t *testing.T) {
	w := watermark.New()
	if w.Behind(nil) {
		t.Fatal("expected false for nil entry")
	}
}

func TestResetClearsHigh(t *testing.T) {
	w := watermark.New()
	_ = w.Advance(entry(time.Now()))
	w.Reset()
	if !w.High().IsZero() {
		t.Fatalf("expected zero after reset, got %v", w.High())
	}
}
