package signal_test

import (
	"sync/atomic"
	"testing"

	"github.com/logpipe/logpipe/internal/signal"
)

func TestMultipleSignalsAreIndependent(t *testing.T) {
	b := signal.New()
	var aCount, bCount int32

	_ = b.Subscribe("a", func(_ string, _ map[string]any) { atomic.AddInt32(&aCount, 1) })
	_ = b.Subscribe("b", func(_ string, _ map[string]any) { atomic.AddInt32(&bCount, 1) })

	b.Fire("a", nil)
	b.Fire("a", nil)
	b.Fire("b", nil)

	if got := atomic.LoadInt32(&aCount); got != 2 {
		t.Fatalf("signal a: expected 2, got %d", got)
	}
	if got := atomic.LoadInt32(&bCount); got != 1 {
		t.Fatalf("signal b: expected 1, got %d", got)
	}
}

func TestResetDoesNotAffectOtherSignals(t *testing.T) {
	b := signal.New()
	var count int32

	_ = b.Subscribe("keep", func(_ string, _ map[string]any) { atomic.AddInt32(&count, 1) })
	_ = b.Subscribe("drop", func(_ string, _ map[string]any) { atomic.AddInt32(&count, 100) })

	b.Reset("drop")
	b.Fire("keep", nil)
	b.Fire("drop", nil)

	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestNameIsPassedToHandler(t *testing.T) {
	b := signal.New()
	var received string

	_ = b.Subscribe("myevent", func(name string, _ map[string]any) {
		received = name
	})
	b.Fire("myevent", nil)

	if received != "myevent" {
		t.Fatalf("expected 'myevent', got %q", received)
	}
}
