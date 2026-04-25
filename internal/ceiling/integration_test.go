package ceiling_test

import (
	"sync"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/ceiling"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(svc string) reader.LogEntry {
	return reader.LogEntry{Service: svc, Message: "hello", Level: "warn"}
}

func TestConcurrentAllowIsSafe(t *testing.T) {
	c, err := ceiling.New(50, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Allow(makeEntry("svc"))
		}()
	}
	wg.Wait()
}

func TestCeilingEnforcesLimitAcrossMultipleServices(t *testing.T) {
	c, err := ceiling.New(2, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	services := []string{"alpha", "beta", "gamma"}
	allowed := make(map[string]int)

	for _, svc := range services {
		for i := 0; i < 5; i++ {
			if c.Allow(makeEntry(svc)) {
				allowed[svc]++
			}
		}
	}

	for _, svc := range services {
		if allowed[svc] != 2 {
			t.Errorf("service %s: expected 2 allowed, got %d", svc, allowed[svc])
		}
	}
}
