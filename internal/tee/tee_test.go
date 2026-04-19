package tee_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/tee"
)

func baseEntry(msg string) reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestNewNoSinksReturnsError(t *testing.T) {
	_, err := tee.New()
	if err == nil {
		t.Fatal("expected error for zero sinks")
	}
}

func TestNewNilSinkReturnsError(t *testing.T) {
	_, err := tee.New(nil)
	if err == nil {
		t.Fatal("expected error for nil sink")
	}
}

func TestApplyForwardsToAllSinks(t *testing.T) {
	var got1, got2 []string

	s1 := func(e reader.LogEntry) { got1 = append(got1, e.Message) }
	s2 := func(e reader.LogEntry) { got2 = append(got2, e.Message) }

	tr, err := tee.New(s1, s2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr.Apply(baseEntry("hello"))
	tr.Apply(baseEntry("world"))

	if len(got1) != 2 || got1[0] != "hello" || got1[1] != "world" {
		t.Errorf("sink1 got %v", got1)
	}
	if len(got2) != 2 || got2[0] != "hello" || got2[1] != "world" {
		t.Errorf("sink2 got %v", got2)
	}
}

func TestApplyDoesNotMutateEntry(t *testing.T) {
	original := baseEntry("immutable")
	var seen reader.LogEntry

	tr, _ := tee.New(func(e reader.LogEntry) {
		e.Message = "mutated"
		seen = e
	})
	tr.Apply(original)

	if original.Message != "immutable" {
		t.Errorf("original entry was mutated")
	}
	_ = seen
}

func TestLenReturnsCorrectCount(t *testing.T) {
	s := func(_ reader.LogEntry) {}
	tr, _ := tee.New(s, s, s)
	if tr.Len() != 3 {
		t.Errorf("expected 3, got %d", tr.Len())
	}
}
