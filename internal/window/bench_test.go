package window

import (
	"testing"
	"time"
)

func BenchmarkAdd(b *testing.B) {
	c := New(10 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Add()
	}
}

func BenchmarkCount(b *testing.B) {
	c := New(10 * time.Second)
	for i := 0; i < 1000; i++ {
		c.Add()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Count()
	}
}

func BenchmarkAddAndCount(b *testing.B) {
	c := New(5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Add()
		c.Count()
	}
}
