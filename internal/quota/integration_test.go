package quota_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/internal/quota"
)

func TestConcurrentAllow(t *testing.T) {
	q, _ := quota.New(quota.Options{MaxEntries: 50, Window: time.Minute})
	var wg sync.WaitGroup
	exceeded := 0
	var mu sync.Mutex
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := q.Allow("svc"); err == quota.ErrQuotaExceeded {
				mu.Lock()
				exceeded++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if exceeded < 50 {
		t.Errorf("expected at least 50 exceeded, got %d", exceeded)
	}
}

func TestMultipleServicesIndependent(t *testing.T) {
	q, _ := quota.New(quota.Options{MaxEntries: 2, Window: time.Minute})
	services := []string{"alpha", "beta", "gamma"}
	for _, svc := range services {
		for i := 0; i < 2; i++ {
			if err := q.Allow(svc); err != nil {
				t.Errorf("service %s attempt %d: unexpected error %v", svc, i, err)
			}
		}
		if err := q.Allow(svc); err != quota.ErrQuotaExceeded {
			t.Errorf("service %s: expected exceeded", svc)
		}
	}
}

// TestWindowExpiry verifies that quota counters reset after the window elapses.
func TestWindowExpiry(t *testing.T) {
	window := 100 * time.Millisecond
	q, _ := quota.New(quota.Options{MaxEntries: 1, Window: window})

	if err := q.Allow("svc"); err != nil {
		t.Fatalf("first allow: unexpected error %v", err)
	}
	if err := q.Allow("svc"); err != quota.ErrQuotaExceeded {
		t.Fatal("expected quota exceeded before window reset")
	}

	time.Sleep(window + 20*time.Millisecond)

	if err := q.Allow("svc"); err != nil {
		t.Errorf("after window expiry: unexpected error %v", err)
	}
}
