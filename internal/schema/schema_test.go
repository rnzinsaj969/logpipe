package schema_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/schema"
)

func TestValidatePassesCompleteEntry(t *testing.T) {
	v := schema.New(schema.Required)
	entry := map[string]any{
		"message": "hello",
		"level":   "info",
		"service": "api",
	}
	if err := v.Validate(entry); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateMissingField(t *testing.T) {
	v := schema.New(schema.Required)
	entry := map[string]any{
		"message": "hello",
		"level":   "info",
	}
	err := v.Validate(entry)
	if err == nil {
		t.Fatal("expected error for missing 'service'")
	}
	if !strings.Contains(err.Error(), "service") {
		t.Errorf("error should mention 'service', got: %v", err)
	}
}

func TestValidateEmptyStringField(t *testing.T) {
	v := schema.New([]string{"message"})
	entry := map[string]any{"message": "   "}
	if err := v.Validate(entry); err == nil {
		t.Fatal("expected error for blank message")
	}
}

func TestValidateMultipleMissingFields(t *testing.T) {
	v := schema.New(schema.Required)
	entry := map[string]any{}
	err := v.Validate(entry)
	if err == nil {
		t.Fatal("expected error")
	}
	for _, f := range schema.Required {
		if !strings.Contains(err.Error(), f) {
			t.Errorf("error should mention %q, got: %v", f, err)
		}
	}
}

func TestNormalizeAddsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	out := schema.Normalize(map[string]any{"message": "hi"})
	ts, ok := out["timestamp"].(string)
	if !ok || ts == "" {
		t.Fatal("expected timestamp string")
	}
	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatalf("invalid timestamp format: %v", err)
	}
	if parsed.Before(before) {
		t.Error("timestamp should not be before test start")
	}
}

func TestNormalizePreservesExistingTimestamp(t *testing.T) {
	out := schema.Normalize(map[string]any{"timestamp": "2024-01-01T00:00:00Z"})
	if out["timestamp"] != "2024-01-01T00:00:00Z" {
		t.Errorf("existing timestamp should not be overwritten")
	}
}

func TestNormalizeLowercasesLevel(t *testing.T) {
	out := schema.Normalize(map[string]any{"level": "  WARN  "})
	if out["level"] != "warn" {
		t.Errorf("expected 'warn', got %v", out["level"])
	}
}

func TestNormalizeDoesNotMutateInput(t *testing.T) {
	input := map[string]any{"level": "ERROR", "message": "boom"}
	schema.Normalize(input)
	if input["level"] != "ERROR" {
		t.Error("input map should not be mutated")
	}
}
