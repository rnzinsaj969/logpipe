package hedge_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/hedge"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(msg string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: "info", Service: "svc"}
}

func TestNewNilPrimaryReturnsError(t *testing.T) {
	_, err := hedge.New(nil, func(reader.LogEntry) error { return nil }, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for nil primary")
	}
}

func TestNewNilFallbackReturnsError(t *testing.T) {
	_, err := hedge.New(func(reader.LogEntry) error { return nil }, nil, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for nil fallback")
	}
}

func TestNewZeroWindowReturnsError(t *testing.T) {
	noop := func(reader.LogEntry) error { return nil }
	_, err := hedge.New(noop, noop, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestPrimarySucceedsNoFallback(t *testing.T) {
	var fallbackCalled atomic.Bool
	primary := func(reader.LogEntry) error { return nil }
	fallback := func(reader.LogEntry) error {
		fallbackCalled.Store(true)
		return nil
	}
	h, err := hedge.New(primary, fallback, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Apply(makeEntry("hello")); err != nil {
		t.Fatalf("unexpected apply error: %v", err)
	}
	time.Sleep(80 * time.Millisecond)
	if fallbackCalled.Load() {
		t.Error("fallback should not have been called when primary succeeds quickly")
	}
}

func TestPrimaryErrorTriggersImmediateFallback(t *testing.T) {
	var fallbackCalled atomic.Bool
	primary := func(reader.LogEntry) error { return errors.New("primary failure") }
	fallback := func(reader.LogEntry) error {
		fallbackCalled.Store(true)
		return nil
	}
	h, err := hedge.New(primary, fallback, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Apply(makeEntry("fail")); err != nil {
		t.Fatalf("unexpected apply error: %v", err)
	}
	if !fallbackCalled.Load() {
		t.Error("fallback should have been called after primary error")
	}
}

func TestSlowPrimaryTriggersHedge(t *testing.T) {
	var fallbackCalled atomic.Bool
	primary := func(reader.LogEntry) error {
		time.Sleep(120 * time.Millisecond)
		return nil
	}
	fallback := func(reader.LogEntry) error {
		fallbackCalled.Store(true)
		return nil
	}
	h, err := hedge.New(primary, fallback, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Apply(makeEntry("slow")); err != nil {
		t.Fatalf("unexpected apply error: %v", err)
	}
	if !fallbackCalled.Load() {
		t.Error("fallback should have been called when primary is slow")
	}
}

func TestBothSinksFailReturnsError(t *testing.T) {
	primary := func(reader.LogEntry) error { return errors.New("primary err") }
	fallback := func(reader.LogEntry) error { return errors.New("fallback err") }
	h, err := hedge.New(primary, fallback, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Apply(makeEntry("both fail")); err == nil {
		t.Error("expected error when both sinks fail")
	}
}
