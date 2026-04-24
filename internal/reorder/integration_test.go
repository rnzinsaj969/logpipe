package reorder_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/reorder"
)

var base = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeEntry(ts time.Time, msg string) reader.LogEntry {
	return reader.LogEntry{Timestamp: ts, Message: msg, Service: "svc", Level: "info"}
}

func TestReorderStreamProducesGlobalOrder(t *testing.T) {
	r, err := reorder.New(reorder.Options{WindowSize: 4, MaxAge: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	inputs := []reader.LogEntry{
		makeEntry(base.Add(3*time.Second), "d"),
		makeEntry(base.Add(1*time.Second), "b"),
		makeEntry(base.Add(2*time.Second), "c"),
		makeEntry(base, "a"),
	}
	var out []reader.LogEntry
	for _, e := range inputs {
		out = append(out, r.Add(e)...)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	expected := []string{"a", "b", "c", "d"}
	for i, e := range out {
		if e.Message != expected[i] {
			t.Errorf("position %d: want %s got %s", i, expected[i], e.Message)
		}
	}
}

func TestReorderDoesNotMutateInput(t *testing.T) {
	r, _ := reorder.New(reorder.Options{WindowSize: 2, MaxAge: time.Second})
	e := makeEntry(base, "original")
	r.Add(e)
	r.Add(makeEntry(base.Add(time.Second), "second"))
	if e.Message != "original" {
		t.Fatal("input entry was mutated")
	}
}
