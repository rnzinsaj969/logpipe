package quota

import (
	"testing"
	"time"
)

func BenchmarkAllowSingleService(b *testing.B) {
	q, _ := New(Options{MaxEntries: b.N + 1, Window: time.Minute})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Allow("svc")
	}
}

func BenchmarkAllowManyServices(b *testing.B) {
	q, _ := New(Options{MaxEntries: 100, Window: time.Minute})
	services := []string{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Allow(services[i%len(services)])
	}
}

func BenchmarkSnapshot(b *testing.B) {
	q, _ := New(Options{MaxEntries: 1000, Window: time.Minute})
	for _, s := range []string{"x", "y", "z"} {
		q.Allow(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Snapshot()
	}
}
