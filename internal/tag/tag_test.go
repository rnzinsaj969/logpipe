package tag_test

import (
	"testing"
	"time"

	"logpipe/internal/reader"
	"logpipe/internal/tag"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "hello",
		Timestamp: time.Unix(0, 0).UTC(),
		Extra:     map[string]any{"existing": "yes"},
	}
}

func TestApplyAddsTags(t *testing.T) {
	tgr, err := tag.New("env", "prod", "region", "us-east-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tgr.Apply(baseEntry())
	if out.Extra["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", out.Extra["env"])
	}
	if out.Extra["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", out.Extra["region"])
	}
}

func TestApplyDoesNotOverwriteExistingExtra(t *testing.T) {
	tgr, err := tag.New("existing", "overridden")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tgr.Apply(baseEntry())
	if out.Extra["existing"] != "yes" {
		t.Errorf("entry Extra should win; got %v", out.Extra["existing"])
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	tgr, _ := tag.New("env", "staging")
	orig := baseEntry()
	_ = tgr.Apply(orig)
	if _, ok := orig.Extra["env"]; ok {
		t.Error("original entry Extra was mutated")
	}
}

func TestNewOddArgumentsReturnsError(t *testing.T) {
	_, err := tag.New("key")
	if err == nil {
		t.Error("expected error for odd number of arguments")
	}
}

func TestNewEmptyKeyReturnsError(t *testing.T) {
	_, err := tag.New("", "value")
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestLen(t *testing.T) {
	tgr, _ := tag.New("a", "1", "b", "2")
	if tgr.Len() != 2 {
		t.Errorf("expected Len=2, got %d", tgr.Len())
	}
}

func TestApplyNoTagsPreservesEntry(t *testing.T) {
	tgr, err := tag.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	orig := baseEntry()
	out := tgr.Apply(orig)
	if out.Message != orig.Message || out.Service != orig.Service {
		t.Error("entry fields changed unexpectedly")
	}
}
