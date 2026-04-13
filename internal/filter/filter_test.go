package filter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  filter.Level
		ok    bool
	}{
		{"debug", filter.LevelDebug, true},
		{"INFO", filter.LevelInfo, true},
		{"Warn", filter.LevelWarn, true},
		{"error", filter.LevelError, true},
		{"unknown", filter.LevelDebug, false},
	}
	for _, tc := range cases {
		got, ok := filter.ParseLevel(tc.input)
		if ok != tc.ok || got != tc.want {
			t.Errorf("ParseLevel(%q) = (%v, %v), want (%v, %v)", tc.input, got, ok, tc.want, tc.ok)
		}
	}
}

func TestCriteriaMatch(t *testing.T) {
	criteria := &filter.Criteria{
		MinLevel: filter.LevelWarn,
		Service:  "api",
		Keyword:  "timeout",
	}

	cases := []struct {
		name  string
		entry filter.LogEntry
		want  bool
	}{
		{"matches all", filter.LogEntry{Service: "api", Level: filter.LevelError, Message: "connection timeout"}, true},
		{"level too low", filter.LogEntry{Service: "api", Level: filter.LevelInfo, Message: "connection timeout"}, false},
		{"wrong service", filter.LogEntry{Service: "worker", Level: filter.LevelError, Message: "connection timeout"}, false},
		{"missing keyword", filter.LogEntry{Service: "api", Level: filter.LevelError, Message: "disk full"}, false},
		{"exact minimum level", filter.LogEntry{Service: "api", Level: filter.LevelWarn, Message: "read timeout"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := criteria.Match(tc.entry); got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCriteriaMatchNoFilters(t *testing.T) {
	c := &filter.Criteria{MinLevel: filter.LevelDebug}
	e := filter.LogEntry{Service: "any", Level: filter.LevelDebug, Message: "hello"}
	if !c.Match(e) {
		t.Error("expected match with empty criteria")
	}
}
