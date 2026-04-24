package extract_test

import (
	"testing"
	"time"

	"logpipe/internal/extract"
	"logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Now(),
		Extra:     map[string]any{"msg_override": "new message", "svc_name": "other"},
	}
}

func TestNewEmptyRulesReturnsError(t *testing.T) {
	_, err := extract.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNewEmptyFromReturnsError(t *testing.T) {
	_, err := extract.New([]extract.Rule{{From: "", To: extract.TargetMessage}})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestNewUnsupportedTargetReturnsError(t *testing.T) {
	_, err := extract.New([]extract.Rule{{From: "key", To: "timestamp"}})
	if err == nil {
		t.Fatal("expected error for unsupported target")
	}
}

func TestExtractPromotesToMessage(t *testing.T) {
	e, err := extract.New([]extract.Rule{{From: "msg_override", To: extract.TargetMessage, Overwrite: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := e.Apply(baseEntry())
	if out.Message != "new message" {
		t.Errorf("expected %q, got %q", "new message", out.Message)
	}
}

func TestExtractNoOverwriteKeepsOriginal(t *testing.T) {
	e, _ := extract.New([]extract.Rule{{From: "msg_override", To: extract.TargetMessage, Overwrite: false}})
	out := e.Apply(baseEntry())
	if out.Message != "hello" {
		t.Errorf("expected original message %q, got %q", "hello", out.Message)
	}
}

func TestExtractPromotesToService(t *testing.T) {
	e, _ := extract.New([]extract.Rule{{From: "svc_name", To: extract.TargetService, Overwrite: true}})
	out := e.Apply(baseEntry())
	if out.Service != "other" {
		t.Errorf("expected service %q, got %q", "other", out.Service)
	}
}

func TestExtractMissingKeyNoChange(t *testing.T) {
	e, _ := extract.New([]extract.Rule{{From: "nonexistent", To: extract.TargetLevel, Overwrite: true}})
	entry := baseEntry()
	out := e.Apply(entry)
	if out.Level != entry.Level {
		t.Errorf("expected level unchanged, got %q", out.Level)
	}
}

func TestExtractDoesNotMutateOriginal(t *testing.T) {
	e, _ := extract.New([]extract.Rule{{From: "msg_override", To: extract.TargetMessage, Overwrite: true}})
	original := baseEntry()
	_ = e.Apply(original)
	if original.Message != "hello" {
		t.Error("original entry was mutated")
	}
}

func TestExtractNonStringValueIgnored(t *testing.T) {
	e, _ := extract.New([]extract.Rule{{From: "count", To: extract.TargetLevel, Overwrite: true}})
	entry := baseEntry()
	entry.Extra["count"] = 42
	out := e.Apply(entry)
	if out.Level != "info" {
		t.Errorf("expected level unchanged, got %q", out.Level)
	}
}
