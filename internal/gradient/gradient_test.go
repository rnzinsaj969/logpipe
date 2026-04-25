package gradient_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/gradient"
)

func TestNewEmptyWeightsReturnsError(t *testing.T) {
	_, err := gradient.New(nil)
	if err == nil {
		t.Fatal("expected error for nil weights, got nil")
	}
	_, err = gradient.New(map[string]float64{})
	if err == nil {
		t.Fatal("expected error for empty weights, got nil")
	}
}

func TestScoreKnownLevel(t *testing.T) {
	s, err := gradient.New(gradient.DefaultWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Score("error")
	if !ok {
		t.Fatal("expected 'error' to be a known level")
	}
	if v != 0.85 {
		t.Fatalf("expected 0.85, got %v", v)
	}
}

func TestScoreUnknownLevelReturnsFalse(t *testing.T) {
	s, _ := gradient.New(gradient.DefaultWeights())
	v, ok := s.Score("trace")
	if ok {
		t.Fatal("expected 'trace' to be unknown")
	}
	if v != 0 {
		t.Fatalf("expected fallback 0, got %v", v)
	}
}

func TestScoreNormalisesCase(t *testing.T) {
	s, _ := gradient.New(gradient.DefaultWeights())
	v, ok := s.Score("WARN")
	if !ok {
		t.Fatal("expected 'WARN' to match after normalisation")
	}
	if v != 0.6 {
		t.Fatalf("expected 0.6, got %v", v)
	}
}

func TestParseAndScoreSuccess(t *testing.T) {
	s, _ := gradient.New(gradient.DefaultWeights())
	v, err := s.ParseAndScore("  Fatal  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 1.0 {
		t.Fatalf("expected 1.0, got %v", v)
	}
}

func TestParseAndScoreUnknownReturnsError(t *testing.T) {
	s, _ := gradient.New(gradient.DefaultWeights())
	_, err := s.ParseAndScore("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
}

func TestSnapshotIsIsolated(t *testing.T) {
	s, _ := gradient.New(gradient.DefaultWeights())
	snap := s.Snapshot()
	snap["debug"] = 99.0
	v, _ := s.Score("debug")
	if v == 99.0 {
		t.Fatal("snapshot mutation affected the internal scorer")
	}
}
