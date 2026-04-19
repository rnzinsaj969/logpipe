package replay_test

import (
	"context"
	"testing"

	"github.com/logpipe/logpipe/internal/replay"
)

func BenchmarkRunHighSpeed(b *testing.B) {
	entries := makeEntries(100)
	for i := 0; i < b.N; i++ {
		r, _ := replay.New(entries, replay.Options{Speed: 1e9})
		out := make(chan interface{ Done() <-chan struct{} }, 200)
		_ = out
		ch := make(chan interface{}, 200)
		_ = ch
		// Use a real channel for the actual benchmark.
		import_out := make(chan interface{}, 200)
		_ = import_out
		_ = r
	}
}

func BenchmarkRunEntries(b *testing.B) {
	entries := makeEntries(50)
	r, _ := replay.New(entries, replay.Options{Speed: 1e9})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := make(chan interface{}, 100)
		_ = out
		_ = r
		_ = context.Background()
	}
}
