// Package fanout provides a Fanout type that distributes a single LogEntry
// to multiple independent sinks. All sinks receive every entry regardless
// of individual sink errors, making it suitable for broadcasting log entries
// to parallel consumers such as file writers, alerting hooks, and metrics
// collectors without one failing sink blocking the others.
package fanout
