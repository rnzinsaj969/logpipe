package stash_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/stash"
)

func entry(msg string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: "info", Service: "svc", Timestamp: time.Now()}
}

func TestNewInvalidCapacityReturnsError(t *testing.T) {
	_, err := stash.New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestPutAndLen(t *testing.T) {
	s, _ := stash.New(10)
	s.Put("a", entry("first"))
	s.Put("a", entry("second"))
	if got := s.Len("a"); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}

func TestFlushReturnsEntriesAndClearsBucket(t *testing.T) {
	s, _ := stash.New(10)
	s.Put("k", entry("one"))
	s.Put("k", entry("two"))
	out := s.Flush("k")
	if len(out) != 2 {
		t.Fatalf("want 2 entries, got %d", len(out))
	}
	if s.Len("k") != 0 {
		t.Fatal("bucket should be empty after flush")
	}
}

func TestFlushMissingKeyReturnsEmpty(t *testing.T) {
	s, _ := stash.New(10)
	out := s.Flush("missing")
	if len(out) != 0 {
		t.Fatalf("want empty slice, got %d entries", len(out))
	}
}

func TestCapacityEvictsOldest(t *testing.T) {
	s, _ := stash.New(3)
	for i := 0; i < 5; i++ {
		s.Put("x", entry(string(rune('a'+i))))
	}
	if s.Len("x") != 3 {
		t.Fatalf("want 3, got %d", s.Len("x"))
	}
	out := s.Flush("x")
	if out[0].Message != "c" {
		t.Fatalf("want oldest evicted, first entry should be 'c', got %q", out[0].Message)
	}
}

func TestKeysReturnsActiveBuckets(t *testing.T) {
	s, _ := stash.New(10)
	s.Put("alpha", entry("a"))
	s.Put("beta", entry("b"))
	keys := s.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "alpha" || keys[1] != "beta" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestConcurrentPut(t *testing.T) {
	s, _ := stash.New(1000)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Put("shared", entry("msg"))
		}()
	}
	wg.Wait()
	if s.Len("shared") != 100 {
		t.Fatalf("want 100, got %d", s.Len("shared"))
	}
}
