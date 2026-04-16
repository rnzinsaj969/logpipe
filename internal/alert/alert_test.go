package alert

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
)

func entry(level, service, message string) reader.LogEntry {
	return reader.LogEntry{
		Level:   level,
		Service: service,
		Message: message,
		Time:    time.Now(),
	}
}

func TestEvaluateMatchesLevel(t *testing.T) {
	a, err := New([]Rule{{Name: "errors", Level: "error"}})
	if err != nil {
		t.Fatal(err)
	}
	a.Evaluate(entry("error", "api", "boom"))
	select {
	case al := <-a.Alerts():
		if al.Rule != "errors" {
			t.Fatalf("expected rule 'errors', got %q", al.Rule)
		}
	default:
		t.Fatal("expected alert, got none")
	}
}

func TestEvaluateNoMatchLevel(t *testing.T) {
	a, _ := New([]Rule{{Name: "errors", Level: "error"}})
	a.Evaluate(entry("info", "api", "ok"))
	select {
	case <-a.Alerts():
		t.Fatal("unexpected alert")
	default:
	}
}

func TestEvaluateMatchesPattern(t *testing.T) {
	a, err := New([]Rule{{Name: "panic", Pattern: `panic`}})
	if err != nil {
		t.Fatal(err)
	}
	a.Evaluate(entry("error", "svc", "panic: runtime error"))
	select {
	case al := <-a.Alerts():
		if al.Rule != "panic" {
			t.Fatalf("expected rule 'panic', got %q", al.Rule)
		}
	default:
		t.Fatal("expected alert")
	}
}

func TestEvaluateCombinedRule(t *testing.T) {
	a, _ := New([]Rule{{Name: "api-error", Level: "error", Service: "api"}})
	a.Evaluate(entry("error", "worker", "fail"))
	select {
	case <-a.Alerts():
		t.Fatal("unexpected alert: service should not match")
	default:
	}
	a.Evaluate(entry("error", "api", "fail"))
	select {
	case al := <-a.Alerts():
		if al.Rule != "api-error" {
			t.Fatalf("unexpected rule %q", al.Rule)
		}
	default:
		t.Fatal("expected alert")
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := New([]Rule{{Name: "bad", Pattern: `[invalid`}})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
