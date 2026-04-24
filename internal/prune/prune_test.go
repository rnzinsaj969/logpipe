package prune_test

import (
	"testing"

	"logpipe/internal/prune"
	"logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service: "svc",
		Level:   "info",
		Message: "hello",
		Extra: map[string]any{
			"token":  "secret",
			"user":   "alice",
			"region": "us-east-1",
		},
	}
}

func TestNewEmptyKeysReturnsError(t *testing.T) {
	_, err := prune.New()
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestNewEmptyStringKeyReturnsError(t *testing.T) {
	_, err := prune.New("")
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestApplyRemovesSingleKey(t *testing.T) {
	p, err := prune.New("token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := p.Apply(base())
	if _, ok := out.Extra["token"]; ok {
		t.Error("expected 'token' to be removed")
	}
	if out.Extra["user"] != "alice" {
		t.Error("expected 'user' to be preserved")
	}
}

func TestApplyRemovesMultipleKeys(t *testing.T) {
	p, _ := prune.New("token", "region")
	out := p.Apply(base())
	if _, ok := out.Extra["token"]; ok {
		t.Error("expected 'token' to be removed")
	}
	if _, ok := out.Extra["region"]; ok {
		t.Error("expected 'region' to be removed")
	}
	if out.Extra["user"] != "alice" {
		t.Error("expected 'user' to be preserved")
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	p, _ := prune.New("token")
	e := base()
	p.Apply(e)
	if _, ok := e.Extra["token"]; !ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApplyNilExtraReturnsUnchanged(t *testing.T) {
	p, _ := prune.New("token")
	e := reader.LogEntry{Service: "svc", Level: "info", Message: "hi"}
	out := p.Apply(e)
	if out.Message != "hi" {
		t.Error("entry should be unchanged")
	}
}

func TestApplyKeyAbsentLeavesEntryUnchanged(t *testing.T) {
	p, _ := prune.New("nonexistent")
	e := base()
	out := p.Apply(e)
	if len(out.Extra) != len(e.Extra) {
		t.Errorf("expected %d keys, got %d", len(e.Extra), len(out.Extra))
	}
}
