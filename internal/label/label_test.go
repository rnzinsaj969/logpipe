package label_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/label"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service: "api",
		Level:   "error",
		Message: "connection refused by host",
		Extra:   map[string]any{"region": "us-east-1"},
	}
}

func TestApplyAddsLabelOnMessageMatch(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "message", Contains: "refused", Label: "alert", Value: "true"},
	})
	out := l.Apply(baseEntry())
	if v, ok := out.Extra["alert"]; !ok || v != "true" {
		t.Fatalf("expected Extra[alert]=true, got %v", out.Extra)
	}
}

func TestApplyAddsLabelOnLevelMatch(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "level", Contains: "error", Label: "severity", Value: "high"},
	})
	out := l.Apply(baseEntry())
	if v, ok := out.Extra["severity"]; !ok || v != "high" {
		t.Fatalf("expected Extra[severity]=high, got %v", out.Extra)
	}
}

func TestApplyAddsLabelOnServiceMatch(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "service", Contains: "api", Label: "team", Value: "platform"},
	})
	out := l.Apply(baseEntry())
	if v, ok := out.Extra["team"]; !ok || v != "platform" {
		t.Fatalf("expected Extra[team]=platform, got %v", out.Extra)
	}
}

func TestApplyAddsLabelOnExtraFieldMatch(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "region", Contains: "us-east", Label: "geo", Value: "americas"},
	})
	out := l.Apply(baseEntry())
	if v, ok := out.Extra["geo"]; !ok || v != "americas" {
		t.Fatalf("expected Extra[geo]=americas, got %v", out.Extra)
	}
}

func TestApplyNoMatchLeavesEntryUnchanged(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "message", Contains: "timeout", Label: "alert", Value: "true"},
	})
	out := l.Apply(baseEntry())
	if _, ok := out.Extra["alert"]; ok {
		t.Fatal("expected no label added for non-matching rule")
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "level", Contains: "error", Label: "x", Value: "1"},
	})
	e := baseEntry()
	_ = l.Apply(e)
	if _, ok := e.Extra["x"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApplyMultipleRulesAllMatch(t *testing.T) {
	l := label.New([]label.Rule{
		{Field: "level", Contains: "error", Label: "sev", Value: "high"},
		{Field: "message", Contains: "refused", Label: "alert", Value: "true"},
	})
	out := l.Apply(baseEntry())
	if out.Extra["sev"] != "high" || out.Extra["alert"] != "true" {
		t.Fatalf("expected both labels set, got %v", out.Extra)
	}
}

func TestHasRulesReturnsFalseWhenEmpty(t *testing.T) {
	l := label.New(nil)
	if l.HasRules() {
		t.Fatal("expected HasRules to return false for empty labeler")
	}
}

func TestHasRulesReturnsTrueWhenPopulated(t *testing.T) {
	l := label.New([]label.Rule{{Field: "level", Contains: "error", Label: "x", Value: "1"}})
	if !l.HasRules() {
		t.Fatal("expected HasRules to return true")
	}
}
