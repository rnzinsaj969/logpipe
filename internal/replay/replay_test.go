package replay_test

import (
	"context"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/replay"
)

func makeEntries(n int) []reader.LogEntry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]reader.LogEntry, n)
	for i := range out {
		out[i] = reader.LogEntry{
			Service:   "svc",
			Level:     "info",
			Message:   "msg",
			Timestamp: base.Add(time.Duration(i) * time.Millisecond),
		}
	}
	return out
}

func TestNewEmptyEntriesReturnsError(t *testing.T) {
	_, err := replay.New(nil, replay.Options{})
	if err == nil {
		t.Fatal("expected error for empty entries")
	}
}

func TestLenReturnsCount(t *testing.T) {
	r, err := replay.New(makeEntries(5), replay.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Len() != 5 {
		t.Fatalf("want 5, got %d", r.Len())
	}
}

func TestRunDeliversAllEntries(t *testing.T) {
	entries := makeEntries(3)
	r, err := replay.New(entries, replay.Options{Speed: 1000.0})
	if err != nil {
		t.Fatal(err)
	}
	out := make(chan reader.LogEntry, 10)
	ctx := context.Background()
	r.Run(ctx, out)
	close(out)
	var got []reader.LogEntry
	for e := range out {
		got = append(got, e)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 entries, got %d", len(got))
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	entries := makeEntries(50)
	r, err := replay.New(entries, replay.Options{Speed: 0.001})
	if err != nil {
		t.Fatal(err)
	}
	out := make(chan reader.LogEntry, 100)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	r.Run(ctx, out)
	if len(out) >= 50 {
		t.Fatal("expected replay to be cancelled before all entries")
	}
}

func TestDefaultSpeedIsOneOnZero(t *testing.T) {
	// speed 0 should default to 1.0 — just ensure no panic
	r, err := replay.New(makeEntries(2), replay.Options{Speed: 0})
	if err != nil {
		t.Fatal(err)
	}
	if r.Len() != 2 {
		t.Fatal("unexpected len")
	}
}
