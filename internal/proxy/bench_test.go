package proxy_test

import (
	"fmt"
	"testing"

	"github.com/logpipe/logpipe/internal/proxy"
)

func BenchmarkForwardSingleSink(b *testing.B) {
	p := proxy.New()
	_ = p.Register("s", &captureSink{})
	e := baseEntry()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Forward(e)
	}
}

func BenchmarkForwardTenSinks(b *testing.B) {
	p := proxy.New()
	for i := 0; i < 10; i++ {
		_ = p.Register(fmt.Sprintf("s%d", i), &captureSink{})
	}
	e := baseEntry()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Forward(e)
	}
}
