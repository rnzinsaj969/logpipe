package gate_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/gate"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Now(),
	}
}

func TestOpenGatePassesEntry(t *testing.T) {
	g, err := gate.New(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := g.Apply(baseEntry())
	if err != nil {
		t.Fatalf("expected entry to pass, got error: %v", err)
	}
	if out.Message != "hello" {
		t.Errorf("unexpected message: %s", out.Message)
	}
}

func TestClosedGateDropsEntry(t *testing.T) {
	g, _ := gate.New(false)
	_, err := g.Apply(baseEntry())
	if err == nil {
		t.Fatal("expected error when gate is closed")
	}
}

func TestOpenAfterClose(t *testing.T) {
	g, _ := gate.New(false)
	g.Open()
	_, err := g.Apply(baseEntry())
	if err != nil {
		t.Fatalf("expected entry to pass after Open(): %v", err)
	}
}

func TestCloseAfterOpen(t *testing.T) {
	g, _ := gate.New(true)
	g.Close()
	_, err := g.Apply(baseEntry())
	if err == nil {
		t.Fatal("expected entry to be dropped after Close()")
	}
}

func TestIsOpen(t *testing.T) {
	g, _ := gate.New(true)
	if !g.IsOpen() {
		t.Error("expected IsOpen to return true")
	}
	g.Close()
	if g.IsOpen() {
		t.Error("expected IsOpen to return false after Close")
	}
}

func TestDoesNotMutateEntry(t *testing.T) {
	g, _ := gate.New(true)
	e := baseEntry()
	out, _ := g.Apply(e)
	if out.Message != e.Message || out.Service != e.Service {
		t.Error("Apply must not mutate the entry")
	}
}
