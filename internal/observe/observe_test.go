package observe_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/observe"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Unix(1_000_000, 0),
	}
}

func TestNewNilHandlerReturnsError(t *testing.T) {
	_, err := observe.New(nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestApplyCallsHandler(t *testing.T) {
	var seen []reader.LogEntry
	obs, err := observe.New(func(e reader.LogEntry) {
		seen = append(seen, e)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := baseEntry()
	out, err := obs.Apply(e)
	if err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}
	if len(seen) != 1 {
		t.Fatalf("expected 1 call, got %d", len(seen))
	}
	if seen[0].Message != e.Message {
		t.Errorf("handler received wrong entry: %+v", seen[0])
	}
	if out.Message != e.Message {
		t.Errorf("Apply mutated entry: %+v", out)
	}
}

func TestApplyDoesNotMutateEntry(t *testing.T) {
	obs, _ := observe.New(func(e reader.LogEntry) {
		e.Message = "mutated"
	})
	e := baseEntry()
	out, _ := obs.Apply(e)
	if out.Message != "hello" {
		t.Errorf("entry was mutated: %+v", out)
	}
}

func TestApplyCalledForEachEntry(t *testing.T) {
	count := 0
	obs, _ := observe.New(func(_ reader.LogEntry) { count++ })

	for i := 0; i < 5; i++ {
		obs.Apply(baseEntry()) //nolint:errcheck
	}
	if count != 5 {
		t.Errorf("expected 5 calls, got %d", count)
	}
}

func TestApplyPreservesAllFields(t *testing.T) {
	obs, _ := observe.New(func(_ reader.LogEntry) {})
	e := baseEntry()
	e.Extra = map[string]any{"key": "value"}

	out, _ := obs.Apply(e)
	if out.Service != e.Service || out.Level != e.Level {
		t.Errorf("fields changed: %+v", out)
	}
	if out.Extra["key"] != "value" {
		t.Errorf("extra field lost: %+v", out.Extra)
	}
}
