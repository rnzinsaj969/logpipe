// Package journal provides a thread-safe, bounded, append-only in-memory
// journal for recording structured log entries.
//
// Each appended entry is assigned a monotonically increasing sequence number
// and tagged with a source identifier, making it easy to correlate entries
// from multiple services during aggregation.
//
// When the journal reaches its configured capacity the oldest entry is
// silently evicted to make room for the newest one (ring behaviour).
//
// Usage:
//
//	j, err := journal.New(1000)
//	if err != nil { ... }
//	j.Append("auth-service", logEntry)
//	entries := j.Snapshot()
package journal
