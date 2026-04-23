package hold_test

import (
	"testing"

	"github.com/logpipe/logpipe/internal/hold"
	"github.com/logpipe/logpipe/internal/reader"
)

func BenchmarkAddNoRelease(b *testing.B) {
	h, _ := hold.New(1024, func(reader.LogEntry) bool { return false })
	e := entry("benchmark message", "info")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(e)
	}
}

func BenchmarkAddWithRelease(b *testing.B) {
	counter := 0
	h, _ := hold.New(64, func(reader.LogEntry) bool {
		counter++
		return counter%64 == 0
	})
	e := entry("benchmark message", "info")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(e)
	}
}
