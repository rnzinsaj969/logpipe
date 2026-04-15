package coalesce_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/coalesce"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Unix(0, 0),
		Extra:     map[string]any{},
	}
}

func TestCoalescePicksFirstNonEmpty(t *testing.T) {
	c, err := coalesce.New([]coalesce.Rule{
		{Sources: []string{"host", "hostname", "node"}, Target: "canonical_host"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := baseEntry()
	e.Extra["hostname"] = "box-2"
	e.Extra["node"] = "box-3"

	out := c.Apply(e)
	if got := out.Extra["canonical_host"]; got != "box-2" {
		t.Errorf("expected box-2, got %v", got)
	}
}

func TestCoalesceSkipsEmptyString(t *testing.T) {
	c, _ := coalesce.New([]coalesce.Rule{
		{Sources: []string{"a", "b"}, Target: "result"},
	})
	e := baseEntry()
	e.Extra["a"] = ""
	e.Extra["b"] = "value-b"

	out := c.Apply(e)
	if got := out.Extra["result"]; got != "value-b" {
		t.Errorf("expected value-b, got %v", got)
	}
}

func TestCoalesceDoesNotMutateOriginal(t *testing.T) {
	c, _ := coalesce.New([]coalesce.Rule{
		{Sources: []string{"src"}, Target: "dst"},
	})
	e := baseEntry()
	e.Extra["src"] = "original"

	_ = c.Apply(e)
	if _, ok := e.Extra["dst"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestCoalesceNoMatchLeavesTargetAbsent(t *testing.T) {
	c, _ := coalesce.New([]coalesce.Rule{
		{Sources: []string{"missing"}, Target: "out"},
	})
	e := baseEntry()
	out := c.Apply(e)
	if _, ok := out.Extra["out"]; ok {
		t.Error("target should not be set when no source matches")
	}
}

func TestNewRejectsEmptyTarget(t *testing.T) {
	_, err := coalesce.New([]coalesce.Rule{{Sources: []string{"a"}, Target: ""}})
	if err == nil {
		t.Error("expected error for empty target")
	}
}

func TestNewRejectsNoSources(t *testing.T) {
	_, err := coalesce.New([]coalesce.Rule{{Sources: []string{}, Target: "dst"}})
	if err == nil {
		t.Error("expected error for empty sources")
	}
}

func TestHasRules(t *testing.T) {
	c, _ := coalesce.New([]coalesce.Rule{{Sources: []string{"a"}, Target: "b"}})
	if !c.HasRules() {
		t.Error("expected HasRules to return true")
	}
}
