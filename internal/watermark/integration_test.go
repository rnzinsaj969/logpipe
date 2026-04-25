package watermark_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/watermark"
)

func TestConcurrentAdvanceIsSafe(t *testing.T) {
	w := watermark.New()
	base := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = w.Advance(entry(base.Add(time.Duration(i) * time.Millisecond)))
		}(i)
	}
	wg.Wait()
	expected := base.Add(49 * time.Millisecond)
	if !w.High().Equal(expected) {
		t.Fatalf("expected high %v, got %v", expected, w.High())
	}
}

func TestWatermarkTracksMonotonicStream(t *testing.T) {
	w := watermark.New()
	base := time.Now().Truncate(time.Second)
	timestamps := []time.Time{
		base,
		base.Add(time.Second),
		base.Add(2 * time.Second),
		base.Add(time.Millisecond), // late / out-of-order
		base.Add(3 * time.Second),
	}
	for _, ts := range timestamps {
		_ = w.Advance(entry(ts))
	}
	expected := base.Add(3 * time.Second)
	if !w.High().Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, w.High())
	}
	if !w.Behind(entry(base.Add(time.Millisecond))) {
		t.Fatal("late entry should be behind watermark")
	}
}
