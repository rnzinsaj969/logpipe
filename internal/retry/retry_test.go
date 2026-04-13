package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func noSleep(_ time.Duration) {}

func newFast(cfg Config) *Retryer {
	r := New(cfg)
	r.sleep = noSleep
	return r
}

func TestSuccessOnFirstAttempt(t *testing.T) {
	r := newFast(Config{MaxAttempts: 3, BaseDelay: time.Millisecond})
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesUpToMax(t *testing.T) {
	r := newFast(Config{MaxAttempts: 4, BaseDelay: time.Millisecond})
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestSuccessAfterRetry(t *testing.T) {
	r := newFast(Config{MaxAttempts: 5, BaseDelay: time.Millisecond})
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retry, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestCancelledContextStopsRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := newFast(Config{MaxAttempts: 10, BaseDelay: time.Millisecond})
	err := r.Do(ctx, func() error { return errTemp })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestBackoffCappedAtMaxDelay(t *testing.T) {
	r := New(Config{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    2 * time.Second,
	})
	if d := r.backoff(10); d > 2*time.Second {
		t.Fatalf("backoff %v exceeds MaxDelay", d)
	}
}

func TestDefaultsApplied(t *testing.T) {
	r := New(Config{})
	if r.cfg.MaxAttempts != 3 {
		t.Fatalf("expected default MaxAttempts=3, got %d", r.cfg.MaxAttempts)
	}
	if r.cfg.BaseDelay != 100*time.Millisecond {
		t.Fatalf("unexpected BaseDelay: %v", r.cfg.BaseDelay)
	}
}
