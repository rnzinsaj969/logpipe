package signal

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSubscribeEmptyNameReturnsError(t *testing.T) {
	b := New()
	if err := b.Subscribe("", func(_ string, _ map[string]any) {}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSubscribeNilHandlerReturnsError(t *testing.T) {
	b := New()
	if err := b.Subscribe("evt", nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestFireCallsSubscriber(t *testing.T) {
	b := New()
	var called int32
	_ = b.Subscribe("ready", func(_ string, _ map[string]any) {
		atomic.AddInt32(&called, 1)
	})
	b.Fire("ready", nil)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestFirePassesPayload(t *testing.T) {
	b := New()
	received := map[string]any{}
	_ = b.Subscribe("data", func(_ string, p map[string]any) {
		for k, v := range p {
			received[k] = v
		}
	})
	b.Fire("data", map[string]any{"key": "value"})
	if received["key"] != "value" {
		t.Fatalf("unexpected payload: %v", received)
	}
}

func TestFireNoSubscribersIsNoop(t *testing.T) {
	b := New()
	b.Fire("unknown", nil) // must not panic
}

func TestMultipleSubscribersAllCalled(t *testing.T) {
	b := New()
	var count int32
	for i := 0; i < 3; i++ {
		_ = b.Subscribe("tick", func(_ string, _ map[string]any) {
			atomic.AddInt32(&count, 1)
		})
	}
	b.Fire("tick", nil)
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestResetRemovesHandlers(t *testing.T) {
	b := New()
	_ = b.Subscribe("ev", func(_ string, _ map[string]any) {})
	b.Reset("ev")
	if b.Len("ev") != 0 {
		t.Fatal("expected 0 handlers after reset")
	}
}

func TestConcurrentFireIsSafe(t *testing.T) {
	b := New()
	_ = b.Subscribe("ping", func(_ string, _ map[string]any) {})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Fire("ping", nil)
		}()
	}
	wg.Wait()
}
