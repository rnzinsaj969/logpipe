package clip_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/clip"
	"github.com/logpipe/logpipe/internal/reader"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Service: "svc",
		Level:   "info",
		Message: "hello",
		Extra:   map[string]any{"score": float64(150), "count": float64(3)},
	}
}

func TestClampAboveMax(t *testing.T) {
	c, err := clip.New([]clip.Rule{{Field: "score", Min: 0, Max: 100}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := c.Apply(baseEntry())
	if got := out.Extra["score"]; got != float64(100) {
		t.Errorf("expected 100, got %v", got)
	}
}

func TestClampBelowMin(t *testing.T) {
	c, _ := clip.New([]clip.Rule{{Field: "score", Min: 200, Max: 300}})
	out := c.Apply(baseEntry())
	if got := out.Extra["score"]; got != float64(200) {
		t.Errorf("expected 200, got %v", got)
	}
}

func TestClampWithinRange(t *testing.T) {
	c, _ := clip.New([]clip.Rule{{Field: "count", Min: 0, Max: 10}})
	out := c.Apply(baseEntry())
	if got := out.Extra["count"]; got != float64(3) {
		t.Errorf("expected 3 unchanged, got %v", got)
	}
}

func TestClampDoesNotMutateOriginal(t *testing.T) {
	c, _ := clip.New([]clip.Rule{{Field: "score", Min: 0, Max: 100}})
	orig := baseEntry()
	c.Apply(orig)
	if orig.Extra["score"] != float64(150) {
		t.Error("original entry was mutated")
	}
}

func TestClampNonNumericFieldUntouched(t *testing.T) {
	c, _ := clip.New([]clip.Rule{{Field: "score", Min: 0, Max: 100}})
	e := baseEntry()
	e.Extra["score"] = "high"
	out := c.Apply(e)
	if got := out.Extra["score"]; got != "high" {
		t.Errorf("expected string 'high', got %v", got)
	}
}

func TestClampMissingFieldIgnored(t *testing.T) {
	c, _ := clip.New([]clip.Rule{{Field: "missing", Min: 0, Max: 1}})
	out := c.Apply(baseEntry())
	if _, ok := out.Extra["missing"]; ok {
		t.Error("unexpected key 'missing' inserted")
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := clip.New([]clip.Rule{{Field: "", Min: 0, Max: 1}})
	if err == nil {
		t.Error("expected error for empty field")
	}
}

func TestNewMinGreaterThanMaxReturnsError(t *testing.T) {
	_, err := clip.New([]clip.Rule{{Field: "score", Min: 10, Max: 5}})
	if err == nil {
		t.Error("expected error for min > max")
	}
}

func TestClampNoRulesReturnsEntryUnchanged(t *testing.T) {
	c, _ := clip.New(nil)
	orig := baseEntry()
	out := c.Apply(orig)
	if out.Extra["score"] != float64(150) {
		t.Error("entry should be unchanged when no rules")
	}
}
