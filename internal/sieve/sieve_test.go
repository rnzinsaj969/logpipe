package sieve_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/reader"
	"github.com/logpipe/logpipe/internal/sieve"
)

func entry(level, service, message string) reader.LogEntry {
	return reader.LogEntry{Level: level, Service: service, Message: message, Timestamp: time.Now()}
}

func TestNoRulesPassesAll(t *testing.T) {
	s, err := sieve.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Apply(entry("error", "api", "boom")) {
		t.Error("expected entry to pass with no rules")
	}
}

func TestMatchByLevel(t *testing.T) {
	s, _ := sieve.New([]sieve.Rule{{Level: "error"}})
	if !s.Apply(entry("error", "api", "msg")) {
		t.Error("expected error entry to pass")
	}
	if s.Apply(entry("info", "api", "msg")) {
		t.Error("expected info entry to be dropped")
	}
}

func TestMatchByService(t *testing.T) {
	s, _ := sieve.New([]sieve.Rule{{Service: "worker"}})
	if !s.Apply(entry("info", "worker", "ok")) {
		t.Error("expected worker entry to pass")
	}
	if s.Apply(entry("info", "api", "ok")) {
		t.Error("expected api entry to be dropped")
	}
}

func TestMatchByPattern(t *testing.T) {
	s, err := sieve.New([]sieve.Rule{{Pattern: "timeout"}})
	if err != nil {
		t.Fatal(err)
	}
	if !s.Apply(entry("warn", "db", "connection timeout")) {
		t.Error("expected pattern match to pass")
	}
	if s.Apply(entry("warn", "db", "all good")) {
		t.Error("expected non-matching entry to be dropped")
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := sieve.New([]sieve.Rule{{Pattern: "[invalid"}})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestMultipleRulesAnyMatch(t *testing.T) {
	s, _ := sieve.New([]sieve.Rule{
		{Level: "error"},
		{Service: "auth"},
	})
	if !s.Apply(entry("info", "auth", "login")) {
		t.Error("expected auth service to pass via second rule")
	}
	if !s.Apply(entry("error", "other", "crash")) {
		t.Error("expected error level to pass via first rule")
	}
	if s.Apply(entry("info", "api", "ping")) {
		t.Error("expected unmatched entry to be dropped")
	}
}
