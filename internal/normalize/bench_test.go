package normalize_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/normalize"
	"github.com/logpipe/logpipe/internal/reader"
)

var sink reader.LogEntry

func BenchmarkApplyDefault(b *testing.B) {
	n := normalize.New(normalize.DefaultOptions())
	e := reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Service:   "  benchmark-svc  ",
		Message:   "  something happened  ",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = n.Apply(e)
	}
}

func BenchmarkApplyNoTrim(b *testing.B) {
	opts := normalize.DefaultOptions()
	opts.TrimSpace = false
	n := normalize.New(opts)
	e := reader.LogEntry{
		Timestamp: time.Now(),
		Level:     "WARN",
		Service:   "svc",
		Message:   "no trim needed",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = n.Apply(e)
	}
}
