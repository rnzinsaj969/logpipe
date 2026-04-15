package format_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/format"
	"github.com/yourorg/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Level:     "info",
		Service:   "api",
		Message:   "request received",
	}
}

func TestFormatTextBasic(t *testing.T) {
	f := format.New(format.DefaultOptions())
	e := baseEntry()
	out := f.FormatText(e)

	if out != "2024-01-15T12:00:00Z [info] api: request received" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTextUpperLevel(t *testing.T) {
	opts := format.DefaultOptions()
	opts.UpperLevel = true
	f := format.New(opts)

	out := f.FormatText(baseEntry())
	if out != "2024-01-15T12:00:00Z [INFO] api: request received" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTextWithExtra(t *testing.T) {
	f := format.New(format.DefaultOptions())
	e := baseEntry()
	e.Extra = map[string]interface{}{"request_id": "abc123"}

	out := f.FormatText(e)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if len(out) <= len("2024-01-15T12:00:00Z [info] api: request received") {
		t.Errorf("expected extra fields in output, got: %q", out)
	}
}

func TestFormatTextCustomTimeFormat(t *testing.T) {
	opts := format.DefaultOptions()
	opts.TimeFormat = "2006-01-02"
	f := format.New(opts)

	out := f.FormatText(baseEntry())
	if out != "2024-01-15 [info] api: request received" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLevel(t *testing.T) {
	opts := format.DefaultOptions()
	opts.UpperLevel = true
	f := format.New(opts)

	if got := f.FormatLevel("warn"); got != "WARN" {
		t.Errorf("expected WARN, got %q", got)
	}
}

func TestFormatLevelNoUpper(t *testing.T) {
	f := format.New(format.DefaultOptions())
	if got := f.FormatLevel("error"); got != "error" {
		t.Errorf("expected error, got %q", got)
	}
}

func TestFormatTimestamp(t *testing.T) {
	opts := format.DefaultOptions()
	opts.TimeFormat = time.Kitchen
	f := format.New(opts)

	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	if got := f.FormatTimestamp(ts); got != "12:00PM" {
		t.Errorf("expected 12:00PM, got %q", got)
	}
}

func TestDefaultOptionsEmptyDelimiterFallback(t *testing.T) {
	opts := format.Options{} // zero value — empty delimiter
	f := format.New(opts)
	e := baseEntry()
	e.Extra = map[string]interface{}{"k": "v"}

	// Should not panic; delimiter defaults to space.
	out := f.FormatText(e)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
}
