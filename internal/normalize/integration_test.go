package normalize_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/normalize"
	"github.com/logpipe/logpipe/internal/reader"
)

func TestNormalizeChainMultipleEntries(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())

	entries := []reader.LogEntry{
		{Timestamp: time.Now(), Level: "ERROR", Service: "  svc  ", Message: "  msg1  "},
		{Timestamp: time.Now(), Level: "", Service: "", Message: "msg2"},
		{Timestamp: time.Now(), Level: "WARN", Service: "gateway", Message: "msg3"},
	}

	expectedLevels := []string{"error", "info", "warn"}
	expectedServices := []string{"svc", "unknown", "gateway"}

	for i, e := range entries {
		out := n.Apply(e)
		if out.Level != expectedLevels[i] {
			t.Errorf("entry %d: expected level %q, got %q", i, expectedLevels[i], out.Level)
		}
		if out.Service != expectedServices[i] {
			t.Errorf("entry %d: expected service %q, got %q", i, expectedServices[i], out.Service)
		}
	}
}

func TestNormalizePreservesTimestampAndExtra(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	now := time.Now()
	e := reader.LogEntry{
		Timestamp: now,
		Level:     "DEBUG",
		Service:   "worker",
		Message:   "task done",
		Extra:     map[string]any{"job_id": "42"},
	}
	out := n.Apply(e)
	if !out.Timestamp.Equal(now) {
		t.Error("timestamp should be preserved")
	}
	if out.Extra["job_id"] != "42" {
		t.Error("extra fields should be preserved")
	}
}
