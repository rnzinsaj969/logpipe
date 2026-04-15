package route

import (
	"testing"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(level, service string) reader.LogEntry {
	return reader.LogEntry{Level: level, Service: service, Message: "msg"}
}

func TestMatchByLevel(t *testing.T) {
	rt, err := New([]Rule{
		{Destination: "errors", Level: "error"},
		{Destination: "default"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := rt.Match(entry("error", "svc")); got != "errors" {
		t.Errorf("expected errors, got %q", got)
	}
	if got := rt.Match(entry("info", "svc")); got != "default" {
		t.Errorf("expected default, got %q", got)
	}
}

func TestMatchByServicePattern(t *testing.T) {
	rt, err := New([]Rule{
		{Destination: "auth-sink", ServicePattern: "^auth"},
		{Destination: "catch-all"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := rt.Match(entry("info", "auth-service")); got != "auth-sink" {
		t.Errorf("expected auth-sink, got %q", got)
	}
	if got := rt.Match(entry("info", "billing")); got != "catch-all" {
		t.Errorf("expected catch-all, got %q", got)
	}
}

func TestMatchNoRuleReturnsEmpty(t *testing.T) {
	rt, err := New([]Rule{
		{Destination: "only-errors", Level: "error"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := rt.Match(entry("info", "svc")); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := New([]Rule{
		{Destination: "d", ServicePattern: "["},
	})
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestMissingDestinationReturnsError(t *testing.T) {
	_, err := New([]Rule{{Level: "info"}})
	if err == nil {
		t.Fatal("expected error for missing destination")
	}
}

func TestRulesReturnsCopy(t *testing.T) {
	original := []Rule{
		{Destination: "a", Level: "warn"},
	}
	rt, _ := New(original)
	copy := rt.Rules()
	copy[0].Destination = "mutated"
	if rt.rules[0].Destination != "a" {
		t.Error("Rules() should return an isolated copy")
	}
}
