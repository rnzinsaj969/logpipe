// Package fork provides a binary routing primitive for log entries.
//
// A Fork evaluates a predicate against each incoming LogEntry and
// dispatches it to one of two sink functions: left when the predicate
// returns true, right otherwise.
//
// This is useful for separating high-severity entries from informational
// ones, or routing entries from a specific service to a dedicated output
// while letting the rest flow through a general pipeline.
//
// Example usage:
//
//	f, err := fork.New(
//		func(e reader.LogEntry) bool { return e.Level == "error" },
//		errorSink,
//		defaultSink,
//	)
package fork
