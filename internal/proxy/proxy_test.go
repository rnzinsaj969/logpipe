package proxy_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/proxy"
	"github.com/logpipe/logpipe/internal/reader"
)

type captureSink struct {
	entries []reader.LogEntry
	err     error
}

func (c *captureSink) Write(e reader.LogEntry) error {
	if c.err != nil {
		return c.err
	}
	c.entries = append(c.entries, e)
	return nil
}

func baseEntry() reader.LogEntry {
	return reader.LogEntry{Service: "svc", Level: "info", Message: "hello", Timestamp: time.Now()}
}

func TestForwardToSingleSink(t *testing.T) {
	p := proxy.New()
	s := &captureSink{}
	_ = p.Register("a", s)
	if err := p.Forward(baseEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.entries))
	}
}

func TestForwardToMultipleSinks(t *testing.T) {
	p := proxy.New()
	a, b := &captureSink{}, &captureSink{}
	_ = p.Register("a", a)
	_ = p.Register("b", b)
	_ = p.Forward(baseEntry())
	if len(a.entries) != 1 || len(b.entries) != 1 {
		t.Fatal("both sinks should receive the entry")
	}
}

func TestRegisterDuplicateReturnsError(t *testing.T) {
	p := proxy.New()
	_ = p.Register("x", &captureSink{})
	if err := p.Register("x", &captureSink{}); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestRegisterEmptyNameReturnsError(t *testing.T) {
	p := proxy.New()
	if err := p.Register("", &captureSink{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRemoveSink(t *testing.T) {
	p := proxy.New()
	_ = p.Register("a", &captureSink{})
	p.Remove("a")
	if p.Len() != 0 {
		t.Fatal("expected 0 sinks after removal")
	}
}

func TestForwardCollectsSinkErrors(t *testing.T) {
	p := proxy.New()
	_ = p.Register("bad", &captureSink{err: errors.New("boom")})
	if err := p.Forward(baseEntry()); err == nil {
		t.Fatal("expected error from failing sink")
	}
}
