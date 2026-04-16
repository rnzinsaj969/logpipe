package field_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/field"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Now(),
		Extra:     map[string]any{"a": "1", "b": "2"},
	}
}

func TestCopyField(t *testing.T) {
	p, err := field.New([]field.Op{{Action: "copy", From: "a", To: "c"}})
	if err != nil {
		t.Fatal(err)
	}
	out := p.Apply(base())
	if out.Extra["a"] != "1" || out.Extra["c"] != "1" {
		t.Fatalf("unexpected extra: %v", out.Extra)
	}
}

func TestRenameField(t *testing.T) {
	p, _ := field.New([]field.Op{{Action: "rename", From: "a", To: "z"}})
	out := p.Apply(base())
	if _, ok := out.Extra["a"]; ok {
		t.Fatal("original key should be removed")
	}
	if out.Extra["z"] != "1" {
		t.Fatalf("expected z=1, got %v", out.Extra["z"])
	}
}

func TestDeleteField(t *testing.T) {
	p, _ := field.New([]field.Op{{Action: "delete", From: "b"}})
	out := p.Apply(base())
	if _, ok := out.Extra["b"]; ok {
		t.Fatal("key b should be deleted")
	}
	if out.Extra["a"] != "1" {
		t.Fatal("key a should remain")
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	p, _ := field.New([]field.Op{{Action: "delete", From: "a"}})
	e := base()
	p.Apply(e)
	if _, ok := e.Extra["a"]; !ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestInvalidOpReturnsError(t *testing.T) {
	_, err := field.New([]field.Op{{Action: "unknown", From: "a"}})
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestMissingFromReturnsError(t *testing.T) {
	_, err := field.New([]field.Op{{Action: "copy", From: "", To: "x"}})
	if err == nil {
		t.Fatal("expected error when From is empty")
	}
}
