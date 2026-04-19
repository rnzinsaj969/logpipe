package proxy_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/logpipe/logpipe/internal/proxy"
)

func TestConcurrentForward(t *testing.T) {
	p := proxy.New()
	var mu sync.Mutex
	var count int
	_ = p.Register("counter", &captureSink{})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = p.Forward(baseEntry())
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	if count != 50 {
		t.Fatalf("expected 50 forwards, got %d", count)
	}
}

func TestRegisterAndRemoveConcurrent(t *testing.T) {
	p := proxy.New()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			name := fmt.Sprintf("sink-%d", n)
			_ = p.Register(name, &captureSink{})
			_ = p.Forward(baseEntry())
			p.Remove(name)
		}(i)
	}
	wg.Wait()
}
