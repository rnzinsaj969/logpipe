package alert_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/alert"
	"github.com/logpipe/logpipe/internal/reader"
)

func makeEntry(level, service, msg string) reader.LogEntry {
	return reader.LogEntry{Level: level, Service: service, Message: msg, Time: time.Now()}
}

func TestMultipleRulesCanFire(t *testing.T) {
	rules := []alert.Rule{
		{Name: "error-rule", Level: "error"},
		{Name: "pattern-rule", Pattern: `timeout`},
	}
	a, err := alert.New(rules)
	if err != nil {
		t.Fatal(err)
	}
	a.Evaluate(makeEntry("error", "svc", "timeout reached"))

	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case al := <-a.Alerts():
			seen[al.Rule] = true
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timed out waiting for alert %d", i)
		}
	}
	if !seen["error-rule"] || !seen["pattern-rule"] {
		t.Fatalf("expected both rules to fire, got %v", seen)
	}
}

func TestNoRulesNoAlerts(t *testing.T) {
	a, _ := alert.New(nil)
	a.Evaluate(makeEntry("error", "svc", "boom"))
	select {
	case <-a.Alerts():
		t.Fatal("unexpected alert with no rules")
	case <-time.After(20 * time.Millisecond):
	}
}
