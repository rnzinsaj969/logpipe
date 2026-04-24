package intern_test

import (
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/intern"
)

func TestInternReturnsSameString(t *testing.T) {
	p := intern.New()
	a := p.Intern("hello")
	b := p.Intern("hello")
	if a != b {
		t.Fatalf("expected same value, got %q and %q", a, b)
	}
}

func TestInternDistinctStrings(t *testing.T) {
	p := intern.New()
	a := p.Intern("info")
	b := p.Intern("error")
	if a == b {
		t.Fatal("distinct strings should not be equal after intern")
	}
}

func TestLenTracksUniqueStrings(t *testing.T) {
	p := intern.New()
	p.Intern("a")
	p.Intern("b")
	p.Intern("a") // duplicate
	if got := p.Len(); got != 2 {
		t.Fatalf("expected Len 2, got %d", got)
	}
}

func TestResetClearsPool(t *testing.T) {
	p := intern.New()
	p.Intern("x")
	p.Reset()
	if got := p.Len(); got != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", got)
	}
}

func TestInternSlice(t *testing.T) {
	p := intern.New()
	ss := []string{"info", "warn", "info"}
	out := p.InternSlice(ss)
	if len(out) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(out))
	}
	if out[0] != out[2] {
		t.Fatal("duplicate entries should be the same interned value")
	}
	if p.Len() != 2 {
		t.Fatalf("expected 2 unique strings, got %d", p.Len())
	}
}

func TestConcurrentIntern(t *testing.T) {
	p := intern.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Intern("shared")
			p.Intern("other")
		}()
	}
	wg.Wait()
	if p.Len() != 2 {
		t.Fatalf("expected 2 unique strings after concurrent access, got %d", p.Len())
	}
}
