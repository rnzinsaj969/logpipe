package transform_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
	"github.com/yourorg/logpipe/internal/transform"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "api",
		Level:     "info",
		Message:   "  hello world  ",
		Timestamp: time.Now(),
	}
}

func TestApplyUpperMessage(t *testing.T) {
	tr := transform.New([]transform.Rule{{Field: "message", Op: "upper"}})
	out := tr.Apply(baseEntry())
	if out.Message != "  HELLO WORLD  " {
		t.Fatalf("expected upper-cased message, got %q", out.Message)
	}
}

func TestApplyLowerLevel(t *testing.T) {
	entry := baseEntry()
	entry.Level = "INFO"
	tr := transform.New([]transform.Rule{{Field: "level", Op: "lower"}})
	out := tr.Apply(entry)
	if out.Level != "info" {
		t.Fatalf("expected lower-cased level, got %q", out.Level)
	}
}

func TestApplyTrimMessage(t *testing.T) {
	tr := transform.New([]transform.Rule{{Field: "message", Op: "trim"}})
	out := tr.Apply(baseEntry())
	if out.Message != "hello world" {
		t.Fatalf("expected trimmed message, got %q", out.Message)
	}
}

func TestApplyPrefixService(t *testing.T) {
	tr := transform.New([]transform.Rule{{Field: "service", Op: "prefix", Value: "svc-"}})
	out := tr.Apply(baseEntry())
	if out.Service != "svc-api" {
		t.Fatalf("expected prefixed service, got %q", out.Service)
	}
}

func TestApplyUnknownOpIsNoop(t *testing.T) {
	entry := baseEntry()
	tr := transform.New([]transform.Rule{{Field: "message", Op: "reverse"}})
	out := tr.Apply(entry)
	if out.Message != entry.Message {
		t.Fatalf("expected no change for unknown op, got %q", out.Message)
	}
}

func TestApplyUnknownFieldIsNoop(t *testing.T) {
	entry := baseEntry()
	tr := transform.New([]transform.Rule{{Field: "metadata", Op: "upper"}})
	out := tr.Apply(entry)
	if out.Message != entry.Message || out.Service != entry.Service {
		t.Fatal("expected no change for unknown field")
	}
}

func TestHasRulesEmpty(t *testing.T) {
	tr := transform.New(nil)
	if tr.HasRules() {
		t.Fatal("expected HasRules to return false for empty transformer")
	}
}

func TestHasRulesNonEmpty(t *testing.T) {
	tr := transform.New([]transform.Rule{{Field: "message", Op: "trim"}})
	if !tr.HasRules() {
		t.Fatal("expected HasRules to return true")
	}
}

func TestApplyMultipleRulesOrdered(t *testing.T) {
	rules := []transform.Rule{
		{Field: "message", Op: "trim"},
		{Field: "message", Op: "upper"},
	}
	tr := transform.New(rules)
	out := tr.Apply(baseEntry())
	if out.Message != "HELLO WORLD" {
		t.Fatalf("expected trimmed then upper-cased message, got %q", out.Message)
	}
}
