package annotate

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func BenchmarkApplySingleRule(b *testing.B) {
	a, _ := New([]Rule{{Pattern: `error`, Key: "flag", Value: "1"}})
	e := reader.LogEntry{
		Service:   "bench",
		Level:     "error",
		Message:   "an error occurred in the service",
		Timestamp: time.Now(),
		Extra:     map[string]any{},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Apply(e)
	}
}

func BenchmarkApplyFiveRules(b *testing.B) {
	a, _ := New([]Rule{
		{Pattern: `error`, Key: "k1", Value: "v1"},
		{Pattern: `timeout`, Key: "k2", Value: "v2"},
		{Pattern: `auth`, Key: "k3", Value: "v3"},
		{Pattern: `disk`, Key: "k4", Value: "v4"},
		{Pattern: `panic`, Key: "k5", Value: "v5"},
	})
	e := reader.LogEntry{
		Service:   "bench",
		Level:     "error",
		Message:   "an error occurred in the service",
		Timestamp: time.Now(),
		Extra:     map[string]any{},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Apply(e)
	}
}
