package reorder

import (
	"testing"
	"time"
)

func BenchmarkAddNoFlush(b *testing.B) {
	r, _ := New(Options{WindowSize: b.N + 1, MaxAge: time.Minute})
	e := entry(base, "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(e)
	}
}

func BenchmarkAddWithFlush(b *testing.B) {
	r, _ := New(Options{WindowSize: 8, MaxAge: time.Minute})
	e := entry(base, "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(e)
	}
}

func BenchmarkFlush(b *testing.B) {
	opts := Options{WindowSize: 64, MaxAge: time.Minute}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r, _ := New(opts)
		for j := 0; j < 64; j++ {
			r.Add(entry(base.Add(time.Duration(64-j)*time.Millisecond), "m"))
		}
		b.StartTimer()
		r.Flush()
	}
}
