package truncate_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/truncate"
)

func baseEntry(msg string) reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Service:   "svc",
		Message:   msg,
		Fields:    map[string]any{},
	}
}

func TestShortMessageUnchanged(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 20, MaxFieldBytes: 10})
	e := baseEntry("hello")
	out := tr.Apply(e)
	if out.Message != "hello" {
		t.Fatalf("expected 'hello', got %q", out.Message)
	}
}

func TestLongMessageTruncated(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 10, MaxFieldBytes: 10})
	long := strings.Repeat("a", 50)
	e := baseEntry(long)
	out := tr.Apply(e)
	if len([]rune(out.Message)) <= 10 {
		// message should end with ellipsis marker
		if !strings.HasSuffix(out.Message, "…") {
			t.Fatalf("expected ellipsis suffix, got %q", out.Message)
		}
	}
	if out.Message == long {
		t.Fatal("expected message to be truncated")
	}
}

func TestFieldStringTruncated(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 100, MaxFieldBytes: 5})
	e := baseEntry("ok")
	e.Fields["token"] = strings.Repeat("x", 20)
	out := tr.Apply(e)
	v, _ := out.Fields["token"].(string)
	if !strings.HasSuffix(v, "…") {
		t.Fatalf("expected truncated field, got %q", v)
	}
}

func TestFieldNonStringUnchanged(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 100, MaxFieldBytes: 5})
	e := baseEntry("ok")
	e.Fields["count"] = 42
	out := tr.Apply(e)
	if out.Fields["count"] != 42 {
		t.Fatalf("expected int field unchanged, got %v", out.Fields["count"])
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 5, MaxFieldBytes: 5})
	orig := strings.Repeat("z", 30)
	e := baseEntry(orig)
	_ = tr.Apply(e)
	if e.Message != orig {
		t.Fatal("Apply must not mutate the original entry")
	}
}

func TestNeededReturnsTrueForLongMessage(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 10, MaxFieldBytes: 10})
	e := baseEntry(strings.Repeat("b", 20))
	if !tr.Needed(e) {
		t.Fatal("expected Needed to return true")
	}
}

func TestNeededReturnsFalseForShortEntry(t *testing.T) {
	tr := truncate.New(truncate.Options{MaxMessageBytes: 100, MaxFieldBytes: 100})
	e := baseEntry("short")
	e.Fields["k"] = "v"
	if tr.Needed(e) {
		t.Fatal("expected Needed to return false")
	}
}

func TestDefaultOptionsUsedOnZeroValues(t *testing.T) {
	tr := truncate.New(truncate.Options{})
	def := truncate.DefaultOptions()
	// A message exactly at the default limit should not be truncated.
	e := baseEntry(strings.Repeat("c", def.MaxMessageBytes))
	out := tr.Apply(e)
	if strings.HasSuffix(out.Message, "…") {
		t.Fatal("message at exact limit should not be truncated")
	}
}
