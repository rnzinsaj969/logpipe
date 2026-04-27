package rank_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/rank"
	"github.com/logpipe/logpipe/internal/reader"
)

func base() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "error",
		Service:   "api",
		Message:   "something failed",
		Extra:     map[string]any{"env": "prod"},
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := rank.New("", []rank.Rule{{Field: "level", Value: "error", Score: 10}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEmptyRulesReturnsError(t *testing.T) {
	_, err := rank.New("priority", nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestApplyScoresByLevel(t *testing.T) {
	r, err := rank.New("priority", []rank.Rule{
		{Field: "level", Value: "error", Score: 100},
		{Field: "level", Value: "warn", Score: 50},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.Apply(base())
	got, ok := out.Extra["priority"]
	if !ok {
		t.Fatal("priority field missing")
	}
	if got.(int) != 100 {
		t.Fatalf("expected 100, got %v", got)
	}
}

func TestApplyAccumulatesMultipleRules(t *testing.T) {
	r, _ := rank.New("score", []rank.Rule{
		{Field: "level", Value: "error", Score: 100},
		{Field: "service", Value: "api", Score: 20},
		{Field: "env", Value: "prod", Score: 5},
	})
	out := r.Apply(base())
	if got := out.Extra["score"].(int); got != 125 {
		t.Fatalf("expected 125, got %d", got)
	}
}

func TestApplyNoMatchReturnsZero(t *testing.T) {
	r, _ := rank.New("score", []rank.Rule{
		{Field: "level", Value: "debug", Score: 10},
	})
	out := r.Apply(base())
	if got := out.Extra["score"].(int); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	r, _ := rank.New("score", []rank.Rule{
		{Field: "level", Value: "error", Score: 10},
	})
	e := base()
	r.Apply(e)
	if _, ok := e.Extra["score"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestApplyExtraFieldRule(t *testing.T) {
	r, _ := rank.New("score", []rank.Rule{
		{Field: "env", Value: "prod", Score: 30},
	})
	out := r.Apply(base())
	if got := out.Extra["score"].(int); got != 30 {
		t.Fatalf("expected 30, got %d", got)
	}
}
