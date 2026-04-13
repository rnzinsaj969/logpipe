package sampling_test

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/sampling"
)

// deterministicSource is a fixed-seed source for reproducible tests.
func deterministicSource(seed int64) rand.Source {
	return rand.NewSource(seed)
}

func TestKeepAlwaysWhenRateOne(t *testing.T) {
	s := sampling.New(1.0, deterministicSource(42))
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatal("expected Keep() == true for rate 1.0")
		}
	}
}

func TestDropAlwaysWhenRateZero(t *testing.T) {
	s := sampling.New(0.0, deterministicSource(42))
	for i := 0; i < 100; i++ {
		if s.Keep() {
			t.Fatal("expected Keep() == false for rate 0.0")
		}
	}
}

func TestRateClampedAboveOne(t *testing.T) {
	s := sampling.New(2.5, deterministicSource(1))
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", s.Rate())
	}
}

func TestRateClampedBelowZero(t *testing.T) {
	s := sampling.New(-0.5, deterministicSource(1))
	if s.Rate() != 0.0 {
		t.Fatalf("expected rate 0.0, got %f", s.Rate())
	}
}

func TestApproximateSamplingRate(t *testing.T) {
	const rate = 0.5
	const iterations = 10_000
	const tolerance = 0.05

	s := sampling.New(rate, rand.NewSource(99))
	kept := 0
	for i := 0; i < iterations; i++ {
		if s.Keep() {
			kept++
		}
	}
	actual := float64(kept) / float64(iterations)
	if actual < rate-tolerance || actual > rate+tolerance {
		t.Fatalf("sampling rate out of tolerance: got %f, want ~%f", actual, rate)
	}
}

func TestSetRateUpdatesRuntime(t *testing.T) {
	s := sampling.New(1.0, deterministicSource(7))
	s.SetRate(0.0)
	if s.Rate() != 0.0 {
		t.Fatalf("expected rate 0.0 after SetRate, got %f", s.Rate())
	}
	if s.Keep() {
		t.Fatal("expected Keep() == false after SetRate(0.0)")
	}
}

func TestConcurrentKeep(t *testing.T) {
	s := sampling.New(0.5, rand.NewSource(123))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s.Keep()
			}
		}()
	}
	wg.Wait()
}
