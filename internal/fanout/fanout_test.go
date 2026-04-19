package fanout_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/fanout"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{Service: "svc", Level: "info", Message: "hello", Timestamp: time.Now()}
}

func TestNewNoSinksReturnsError(t *testing.T) {
	_, err := fanout.New()
	if err == nil {
		t.Fatal("expected error for zero sinks")
	}
}

func TestNewNilSinkReturnsError(t *testing.T) {
	_, err := fanout.New(nil)
	if err == nil {
		t.Fatal("expected error for nil sink")
	}
}

func TestApplyForwardsToAllSinks(t *testing.T) {
	counts := make([]int, 3)
	makeSink := func(i int) fanout.Sink {
		return func(e reader.LogEntry) error { counts[i]++; return nil }
	}
	f, err := fanout.New(makeSink(0), makeSink(1), makeSink(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.Apply(base()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, c := range counts {
		if c != 1 {
			t.Errorf("sink %d: expected 1 call, got %d", i, c)
		}
	}
}

func TestApplyCallsAllSinksEvenOnError(t *testing.T) {
	called := 0
	sinkErr := errors.New("sink error")
	errSink := func(e reader.LogEntry) error { called++; return sinkErr }
	okSink := func(e reader.LogEntry) error { called++; return nil }

	f, _ := fanout.New(errSink, okSink, errSink)
	err := f.Apply(base())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if called != 3 {
		t.Errorf("expected 3 calls, got %d", called)
	}
}

func TestLenReturnsCount(t *testing.T) {
	noop := func(e reader.LogEntry) error { return nil }
	f, _ := fanout.New(noop, noop)
	if f.Len() != 2 {
		t.Errorf("expected 2, got %d", f.Len())
	}
}
