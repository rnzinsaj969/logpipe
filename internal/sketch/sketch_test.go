package sketch

import (
	"sync"
	"testing"
)

func TestNewZeroWidthReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero width, got nil")
	}
}

func TestAddAndCountSingleKey(t *testing.T) {
	s, err := New(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.Add("api", 1)
	s.Add("api", 1)
	s.Add("api", 1)
	if got := s.Count("api"); got < 3 {
		t.Fatalf("expected count >= 3, got %d", got)
	}
}

func TestCountUnseenKeyIsZero(t *testing.T) {
	s, _ := New(64)
	if got := s.Count("unknown"); got != 0 {
		t.Fatalf("expected 0 for unseen key, got %d", got)
	}
}

func TestAddZeroDeltaIgnored(t *testing.T) {
	s, _ := New(64)
	s.Add("key", 0)
	if got := s.Count("key"); got != 0 {
		t.Fatalf("expected 0 after zero-delta add, got %d", got)
	}
}

func TestResetClearsAllCounts(t *testing.T) {
	s, _ := New(64)
	s.Add("svc-a", 10)
	s.Add("svc-b", 5)
	s.Reset()
	if got := s.Count("svc-a"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := s.Count("svc-b"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestMultipleDistinctKeys(t *testing.T) {
	s, _ := New(256)
	keys := []string{"auth", "billing", "gateway", "worker"}
	for i, k := range keys {
		s.Add(k, uint64(i+1))
	}
	for i, k := range keys {
		want := uint64(i + 1)
		if got := s.Count(k); got < want {
			t.Errorf("key %q: expected count >= %d, got %d", k, want, got)
		}
	}
}

func TestConcurrentAdd(t *testing.T) {
	s, _ := New(128)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Add("concurrent", 1)
		}()
	}
	wg.Wait()
	if got := s.Count("concurrent"); got < 100 {
		t.Fatalf("expected count >= 100, got %d", got)
	}
}
