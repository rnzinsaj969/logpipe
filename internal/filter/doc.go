// Package filter provides log entry filtering primitives for logpipe.
//
// It defines the Criteria type, which encapsulates filtering rules such as
// minimum severity level, service name, and keyword matching. Log entries
// represented by LogEntry can be tested against a Criteria using Match.
//
// Example usage:
//
//	minLevel, _ := filter.ParseLevel("warn")
//	c := &filter.Criteria{
//		MinLevel: minLevel,
//		Service:  "api",
//		Keyword:  "timeout",
//	}
//
//	entry := filter.LogEntry{
//		Service: "api",
//		Level:   filter.LevelError,
//		Message: "upstream timeout",
//	}
//
//	if c.Match(entry) {
//		fmt.Println("entry passed filter")
//	}
package filter
