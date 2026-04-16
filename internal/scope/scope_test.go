package scope_test

import (
	"testing"
	"time"

	"logpipe/internal/reader"
	"logpipe/internal/scope"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Time{},
		Extra: map[string]any{
			"meta": map[string]any{
				"region": "us-east",
				"host":   "node1",
			},
		},
	}
}

func TestNewEmptyNamespaceReturnsError(t *testing.T) {
	_, err := scope.New("")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestExtractReturnsNestedMap(t *testing.T) {
	s, _ := scope.New("meta")
	fields := s.Extract(base())
	if fields["region"] != "us-east" {
		t.Errorf("expected us-east, got %v", fields["region"])
	}
	if fields["host"] != "node1" {
		t.Errorf("expected node1, got %v", fields["host"])
	}
}

func TestExtractMissingKeyReturnsEmpty(t *testing.T) {
	s, _ := scope.New("absent")
	fields := s.Extract(base())
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestExtractNonMapValueReturnsEmpty(t *testing.T) {
	s, _ := scope.New("meta")
	entry := base()
	entry.Extra["meta"] = "not-a-map"
	fields := s.Extract(entry)
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestEmbedSetsNamespace(t *testing.T) {
	s, _ := scope.New("meta")
	newFields := map[string]any{"env": "prod"}
	out := s.Embed(base(), newFields)
	v, ok := out.Extra["meta"]
	if !ok {
		t.Fatal("expected meta key in extra")
	}
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatal("expected map under meta")
	}
	if m["env"] != "prod" {
		t.Errorf("expected prod, got %v", m["env"])
	}
}

func TestEmbedDoesNotMutateOriginal(t *testing.T) {
	s, _ := scope.New("meta")
	original := base()
	s.Embed(original, map[string]any{"x": 1})
	m := original.Extra["meta"].(map[string]any)
	if _, ok := m["x"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestKeyReturnsNamespace(t *testing.T) {
	s, _ := scope.New("ctx")
	if s.Key() != "ctx" {
		t.Errorf("expected ctx, got %s", s.Key())
	}
}
