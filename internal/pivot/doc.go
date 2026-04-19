// Package pivot provides a thread-safe aggregation table that groups
// log entries by a chosen field (level, service, or any Extra key) and
// counts occurrences. Call Snapshot to retrieve and reset the current
// counts atomically.
package pivot
