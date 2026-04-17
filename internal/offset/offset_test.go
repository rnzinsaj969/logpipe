package offset_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/offset"
)

func TestGetReturnsZeroForUnknownSource(t *testing.T) {
	s := offset.New()
	if got := s.Get("svc"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSetAndGet(t *testing.T) {
	s := offset.New()
	s.Set("svc", 42)
	if got := s.Get("svc"); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

func TestSetIgnoresNegativeOffset(t *testing.T) {
	s := offset.New()
	s.Set("svc", -1)
	if got := s.Get("svc"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestDelete(t *testing.T) {
	s := offset.New()
	s.Set("svc", 10)
	s.Delete("svc")
	if got := s.Get("svc"); got != 0 {
		t.Fatalf("expected 0 after delete, got %d", got)
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	s := offset.New()
	s.Set("a", 1)
	snap := s.Snapshot()
	snap["a"] = 999
	if got := s.Get("a"); got != 1 {
		t.Fatalf("snapshot mutation affected store: got %d", got)
	}
}

func TestReset(t *testing.T) {
	s := offset.New()
	s.Set("a", 5)
	s.Set("b", 10)
	s.Reset()
	if snap := s.Snapshot(); len(snap) != 0 {
		t.Fatalf("expected empty snapshot after reset, got %v", snap)
	}
}

func TestConcurrentSetAndGet(t *testing.T) {
	s := offset.New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int64) {
			defer wg.Done()
			s.Set("svc", n)
			_ = s.Get("svc")
		}(int64(i))
	}
	wg.Wait()
}
