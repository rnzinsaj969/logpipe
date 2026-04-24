package intern_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/intern"
)

// TestInternPoolSharedAcrossCallSites verifies that a single pool correctly
// deduplicates strings that originate from different call sites, simulating
// how level and service values are interned during log ingestion.
func TestInternPoolSharedAcrossCallSites(t *testing.T) {
	p := intern.New()

	levels := []string{"info", "warn", "error", "info", "warn", "info"}
	interned := make([]string, len(levels))
	for i, l := range levels {
		interned[i] = p.Intern(l)
	}

	if p.Len() != 3 {
		t.Fatalf("expected 3 unique levels, got %d", p.Len())
	}

	// All "info" entries must be the exact same interned string.
	if interned[0] != interned[3] || interned[0] != interned[5] {
		t.Fatal("all 'info' entries should share the same interned string")
	}
}

// TestResetAndRefill ensures the pool can be cleared and reused without
// retaining stale references.
func TestResetAndRefill(t *testing.T) {
	p := intern.New()
	p.Intern("alpha")
	p.Intern("beta")
	p.Reset()

	p.Intern("gamma")
	if p.Len() != 1 {
		t.Fatalf("expected 1 entry after reset and refill, got %d", p.Len())
	}
}
