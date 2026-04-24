package annotate

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "user logged in",
		Timestamp: time.Now(),
		Extra:     map[string]any{},
	}
}

func TestNewEmptyRulesReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNewEmptyKeyReturnsError(t *testing.T) {
	_, err := New([]Rule{{Pattern: ".*", Key: ""}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewInvalidPatternReturnsError(t *testing.T) {
	_, err := New([]Rule{{Pattern: "[", Key: "tag"}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApplyAddsAnnotationOnMatch(t *testing.T) {
	a, err := New([]Rule{{Pattern: "logged in", Key: "auth", Value: "true"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := a.Apply(baseEntry())
	if out.Extra["auth"] != "true" {
		t.Errorf("expected auth=true, got %v", out.Extra["auth"])
	}
}

func TestApplyNoMatchLeavesExtraUnchanged(t *testing.T) {
	a, err := New([]Rule{{Pattern: "panic", Key: "severity", Value: "critical"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := a.Apply(baseEntry())
	if _, ok := out.Extra["severity"]; ok {
		t.Error("expected no severity annotation")
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	a, err := New([]Rule{{Pattern: "logged", Key: "tag", Value: "login"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := baseEntry()
	a.Apply(e)
	if _, ok := e.Extra["tag"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestApplyMultipleRulesAllMatch(t *testing.T) {
	a, err := New([]Rule{
		{Pattern: "logged", Key: "event", Value: "auth"},
		{Pattern: "user", Key: "actor", Value: "human"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := a.Apply(baseEntry())
	if out.Extra["event"] != "auth" {
		t.Errorf("expected event=auth, got %v", out.Extra["event"])
	}
	if out.Extra["actor"] != "human" {
		t.Errorf("expected actor=human, got %v", out.Extra["actor"])
	}
}
