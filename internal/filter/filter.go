package filter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// ParseLevel converts a string to a Level, returning LevelDebug and false if unknown.
func ParseLevel(s string) (Level, bool) {
	l, ok := levelNames[strings.ToLower(s)]
	return l, ok
}

// Criteria holds the filtering parameters applied to log entries.
type Criteria struct {
	MinLevel  Level
	Service   string // empty means all services
	Keyword   string // empty means no keyword filter
}

// LogEntry represents a single structured log line.
type LogEntry struct {
	Service string
	Level   Level
	Message string
}

// Match reports whether the entry satisfies all criteria.
func (c *Criteria) Match(e LogEntry) bool {
	if e.Level < c.MinLevel {
		return false
	}
	if c.Service != "" && !strings.EqualFold(e.Service, c.Service) {
		return false
	}
	if c.Keyword != "" && !strings.Contains(e.Message, c.Keyword) {
		return false
	}
	return true
}
