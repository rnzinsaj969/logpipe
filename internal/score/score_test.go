package score_test

import (
	"testing"
	"time"

	"logpipe/internal/reader"
	"logpipe/internal/score"
)

func baseEntry() reader.LogEntry {
	return reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "error",
		Service:   "auth",
		Message:   "something failed",
	}
}

func TestNewEmptyRulesReturnsError(t *testing.T) {
	_, err := score.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := score.New([]score.Rule{{Field: "", Value: "error", Weight: 10}})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestScoreMatchesSingleRule(t *testing.T) {
	s, err := score.New([]score.Rule{
		{Field: "level", Value: "error", Weight: 10},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Score(baseEntry())
	if got != 10 {
		t.Fatalf("expected 10, got %v", got)
	}
}

func TestScoreNoMatchReturnsZero(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "level", Value: "debug", Weight: 5},
	})
	got := s.Score(baseEntry())
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestScoreAccumulatesMultipleRules(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "level", Value: "error", Weight: 10},
		{Field: "service", Value: "auth", Weight: 5},
	})
	got := s.Score(baseEntry())
	if got != 15 {
		t.Fatalf("expected 15, got %v", got)
	}
}

func TestApplySetsScoreInExtra(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "level", Value: "error", Weight: 7},
	})
	out := s.Apply(baseEntry())
	v, ok := out.Extra["_score"]
	if !ok {
		t.Fatal("expected _score key in Extra")
	}
	if v.(float64) != 7 {
		t.Fatalf("expected 7, got %v", v)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "level", Value: "error", Weight: 3},
	})
	e := baseEntry()
	e.Extra = map[string]any{"env": "prod"}
	_ = s.Apply(e)
	if _, ok := e.Extra["_score"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestScoreMatchesExtraField(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "env", Value: "prod", Weight: 20},
	})
	e := baseEntry()
	e.Extra = map[string]any{"env": "prod"}
	got := s.Score(e)
	if got != 20 {
		t.Fatalf("expected 20, got %v", got)
	}
}

func TestScoreCaseInsensitive(t *testing.T) {
	s, _ := score.New([]score.Rule{
		{Field: "level", Value: "ERROR", Weight: 8},
	})
	got := s.Score(baseEntry()) // level is "error" (lowercase)
	if got != 8 {
		t.Fatalf("expected 8, got %v", got)
	}
}
