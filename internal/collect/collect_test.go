package collect_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/collect"
	"github.com/logpipe/logpipe/internal/reader"
)

func entry(msg string) reader.LogEntry {
	return reader.LogEntry{Message: msg, Level: "info", Service: "svc", Timestamp: time.Now()}
}

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := collect.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestAddAndLen(t *testing.T) {
	c, _ := collect.New(10)
	c.Add(entry("a"))
	c.Add(entry("b"))
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestFlushReturnsEntriesAndResets(t *testing.T) {
	c, _ := collect.New(10)
	c.Add(entry("x"))
	c.Add(entry("y"))
	out := c.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if c.Len() != 0 {
		t.Fatal("expected buffer to be empty after flush")
	}
}

func TestFlushIsIsolated(t *testing.T) {
	c, _ := collect.New(10)
	c.Add(entry("z"))
	out := c.Flush()
	out[0].Message = "mutated"
	c.Add(entry("z"))
	out2 := c.Flush()
	if out2[0].Message == "mutated" {
		t.Fatal("flush snapshot shares memory with internal buffer")
	}
}

func TestRingEvictsOldest(t *testing.T) {
	c, _ := collect.New(3)
	for i, m := range []string{"a", "b", "c", "d"} {
		_ = i
		c.Add(entry(m))
	}
	out := c.Flush()
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
	if out[0].Message != "b" {
		t.Fatalf("expected oldest to be evicted, got %q", out[0].Message)
	}
}

func TestConcurrentAdd(t *testing.T) {
	c, _ := collect.New(100)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				c.Add(entry("concurrent"))
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if c.Len() > 100 {
		t.Fatal("collector exceeded max capacity")
	}
}
