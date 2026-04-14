package checkpoint_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/logpipe/internal/checkpoint"
)

func BenchmarkSet(b *testing.B) {
	s, _ := checkpoint.New(tempPath(b))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(fmt.Sprintf("svc-%d", i%10), checkpoint.State{Offset: int64(i)})
	}
}

func BenchmarkFlush(b *testing.B) {
	path := tempPath(b)
	s, _ := checkpoint.New(path)
	for i := 0; i < 10; i++ {
		s.Set(fmt.Sprintf("svc-%d", i), checkpoint.State{Offset: int64(i * 100), Inode: uint64(i + 1)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Flush()
	}
}

func tempPath(b interface{ TempDir() string }) string {
	return b.TempDir() + "/checkpoint.json"
}
