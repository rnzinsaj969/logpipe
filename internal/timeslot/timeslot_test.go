package timeslot_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
	"github.com/yourorg/logpipe/internal/timeslot"
)

func entry(ts time.Time, msg string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Timestamp: ts}
}

func TestInvalidSizeReturnsError(t *testing.T) {
	_, err := timeslot.New(timeslot.Options{Size: 0})
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestAddGroupsIntoSlots(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	b, _ := timeslot.New(timeslot.Options{Size: time.Minute})

	b.Add(entry(now, "a"))
	b.Add(entry(now.Add(30*time.Second), "b"))
	b.Add(entry(now.Add(90*time.Second), "c"))

	slots := b.Snapshot()
	if len(slots) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(slots))
	}
}

func TestSnapshotResetsState(t *testing.T) {
	now := time.Now()
	b, _ := timeslot.New(timeslot.Options{Size: time.Minute})
	b.Add(entry(now, "x"))
	b.Snapshot()
	slots := b.Snapshot()
	if len(slots) != 0 {
		t.Fatalf("expected empty snapshot after reset, got %d", len(slots))
	}
}

func TestZeroTimestampUsesClockTime(t *testing.T) {
	fixed := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	b, _ := timeslot.New(timeslot.Options{
		Size:  time.Hour,
		Clock: func() time.Time { return fixed },
	})
	b.Add(reader.LogEntry{Message: "no-ts"})
	slots := b.Snapshot()
	if len(slots) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(slots))
	}
	if !slots[0].Start.Equal(fixed.Truncate(time.Hour)) {
		t.Errorf("unexpected slot start: %v", slots[0].Start)
	}
}

func TestMultipleEntriesSameSlot(t *testing.T) {
	now := time.Now().Truncate(time.Minute)
	b, _ := timeslot.New(timeslot.Options{Size: time.Minute})
	for i := 0; i < 5; i++ {
		b.Add(entry(now, "msg"))
	}
	slots := b.Snapshot()
	if len(slots) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(slots))
	}
	if len(slots[0].Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(slots[0].Entries))
	}
}
