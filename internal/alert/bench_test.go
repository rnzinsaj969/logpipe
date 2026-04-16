package alert

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

var sink Alert

func BenchmarkEvaluateNoMatch(b *testing.B) {
	a, _ := New([]Rule{
		{Name: "r1", Level: "error"},
		{Name: "r2", Pattern: `panic`},
	})
	e := reader.LogEntry{Level: "info", Service: "svc", Message: "all good", Time: time.Now()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Evaluate(e)
	}
}

func BenchmarkEvaluateMatch(b *testing.B) {
	a, _ := New([]Rule{
		{Name: "r1", Level: "error"},
		{Name: "r2", Pattern: `panic`},
	})
	e := reader.LogEntry{Level: "error", Service: "api", Message: "panic: nil ptr", Time: time.Now()}
	// drain goroutine
	go func() {
		for range a.Alerts() {
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Evaluate(e)
	}
}
