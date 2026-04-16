package pin_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/pin"
	"github.com/logpipe/logpipe/internal/reader"
)

func BenchmarkApplySingleRule(b *testing.B) {
	p, _ := pin.New([]pin.Rule{{Key: "lvl", Target: "level"}})
	entry := reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "msg",
		Timestamp: time.Now(),
		Extra:     map[string]any{"lvl": "warn", "x": "y"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Apply(entry)
	}
}

func BenchmarkApplyMultipleRules(b *testing.B) {
	p, _ := pin.New([]pin.Rule{
		{Key: "msg_override", Target: "message"},
		{Key: "lvl", Target: "level"},
		{Key: "src", Target: "service"},
	})
	entry := reader.LogEntry{
		Service:   "svc",
		Level:     "info",
		Message:   "msg",
		Timestamp: time.Now(),
		Extra: map[string]any{
			"msg_override": "new",
			"lvl":          "error",
			"src":          "auth",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Apply(entry)
	}
}
