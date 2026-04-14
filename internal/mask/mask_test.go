package mask_test

import (
	"testing"
	"time"

	"github.com/your-org/logpipe/internal/mask"
	"github.com/your-org/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Service:   "auth",
		Level:     "info",
		Message:   "user logged in",
		Extra:     map[string]interface{}{"token": "abc123", "user": "alice"},
	}
}

func TestMaskMessage(t *testing.T) {
	m, err := mask.New(mask.Options{Fields: []string{"message"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(baseEntry())
	if out.Message != "***" {
		t.Errorf("expected message masked, got %q", out.Message)
	}
	if out.Service != "auth" {
		t.Errorf("service should be unchanged, got %q", out.Service)
	}
}

func TestMaskService(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"service"}})
	out := m.Apply(baseEntry())
	if out.Service != "***" {
		t.Errorf("expected service masked, got %q", out.Service)
	}
	if out.Message != "user logged in" {
		t.Errorf("message should be unchanged")
	}
}

func TestMaskExtraField(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"token"}})
	out := m.Apply(baseEntry())
	if out.Extra["token"] != "***" {
		t.Errorf("expected token masked, got %v", out.Extra["token"])
	}
	if out.Extra["user"] != "alice" {
		t.Errorf("user should be unchanged, got %v", out.Extra["user"])
	}
}

func TestMaskDoesNotMutateOriginal(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"message", "token"}})
	orig := baseEntry()
	_ = m.Apply(orig)
	if orig.Message != "user logged in" {
		t.Errorf("original message mutated")
	}
	if orig.Extra["token"] != "abc123" {
		t.Errorf("original extra mutated")
	}
}

func TestCustomPlaceholder(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"message"}, Placeholder: "[REDACTED]"})
	out := m.Apply(baseEntry())
	if out.Message != "[REDACTED]" {
		t.Errorf("expected custom placeholder, got %q", out.Message)
	}
}

func TestNoFieldsReturnsError(t *testing.T) {
	_, err := mask.New(mask.Options{})
	if err == nil {
		t.Error("expected error for empty fields, got nil")
	}
}

func TestHasField(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"token", "message"}})
	if !m.HasField("token") {
		t.Error("expected HasField true for token")
	}
	if m.HasField("service") {
		t.Error("expected HasField false for service")
	}
}

func TestMaskNilExtra(t *testing.T) {
	m, _ := mask.New(mask.Options{Fields: []string{"token"}})
	entry := baseEntry()
	entry.Extra = nil
	out := m.Apply(entry) // must not panic
	if out.Extra != nil {
		t.Errorf("expected nil extra, got %v", out.Extra)
	}
}
