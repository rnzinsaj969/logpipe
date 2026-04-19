package pivot_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/pivot"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(level, service string) reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   service,
		Message:   "test",
	}
}

func TestConcurrentAdd(t *testing.T) {
	tbl, err := pivot.New("level")
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tbl.Add(makeEntry("info", "svc"))
		}()
	}
	wg.Wait()
	snap := tbl.Snapshot()
	if snap["info"] != 100 {
		t.Errorf("expected 100, got %d", snap["info"])
	}
}

func TestMultipleSnapshotsAreIndependent(t *testing.T) {
	tbl, _ := pivot.New("service")
	tbl.Add(makeEntry("info", "alpha"))
	snap1 := tbl.Snapshot()
	tbl.Add(makeEntry("info", "beta"))
	snap2 := tbl.Snapshot()

	if _, ok := snap1["beta"]; ok {
		t.Error("snap1 should not contain beta")
	}
	if _, ok := snap2["alpha"]; ok {
		t.Error("snap2 should not contain alpha")
	}
	if snap2["beta"] != 1 {
		t.Errorf("expected beta=1 in snap2, got %d", snap2["beta"])
	}
}
