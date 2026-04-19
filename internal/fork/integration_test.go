package fork_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/fork"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(level, svc string) reader.LogEntry {
	return reader.LogEntry{Timestamp: time.Now(), Level: level, Service: svc, Message: "m"}
}

func TestForkSplitsStream(t *testing.T) {
	var mu sync.Mutex
	var errors, others []reader.LogEntry

	f, err := fork.New(
		func(e reader.LogEntry) bool { return e.Level == "error" },
		func(e reader.LogEntry) error { mu.Lock(); errors = append(errors, e); mu.Unlock(); return nil },
		func(e reader.LogEntry) error { mu.Lock(); others = append(others, e); mu.Unlock(); return nil },
	)
	if err != nil {
		t.Fatal(err)
	}

	entries := []reader.LogEntry{
		makeEntry("error", "a"),
		makeEntry("info", "b"),
		makeEntry("error", "c"),
		makeEntry("debug", "d"),
	}
	for _, e := range entries {
		if err := f.Apply(e); err != nil {
			t.Fatal(err)
		}
	}
	if len(errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errors))
	}
	if len(others) != 2 {
		t.Fatalf("expected 2 others, got %d", len(others))
	}
}
