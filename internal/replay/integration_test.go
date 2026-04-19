package replay_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/replay"
)

func TestReplayOrderPreserved(t *testing.T) {
	entries := makeEntries(5)
	r, _ := replay.New(entries, replay.Options{Speed: 10000.0})
	out := make(chan reader.LogEntry, 10)
	r.Run(context.Background(), out)
	close(out)
	prev := time.Time{}
	for e := range out {
		if !e.Timestamp.After(prev) && !e.Timestamp.Equal(prev) {
			t.Fatalf("out-of-order entry: %v after %v", e.Timestamp, prev)
		}
		prev = e.Timestamp
	}
}

func TestReplayDoesNotMutateInput(t *testing.T) {
	original := makeEntries(3)
	copy := make([]reader.LogEntry, len(original))
	for i, e := range original {
		copy[i] = e
	}
	r, _ := replay.New(original, replay.Options{Speed: 10000.0})
	out := make(chan reader.LogEntry, 10)
	r.Run(context.Background(), out)
	for i, e := range original {
		if e != copy[i] {
			t.Fatalf("entry %d mutated", i)
		}
	}
}
