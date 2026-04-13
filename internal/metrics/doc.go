// Package metrics provides lightweight, thread-safe counters for tracking
// logpipe pipeline activity at runtime.
//
// Metrics are collected throughout the pipeline — from log entry ingestion
// through filtering and output — and can be queried via Snapshot for
// reporting or diagnostics without blocking the hot path.
//
// Usage:
//
//	m := metrics.New()
//	m.EntriesRead.Inc()
//	m.EntriesMatched.Inc()
//	snap := m.Snapshot()
//	fmt.Printf("read=%d matched=%d\n", snap.EntriesRead, snap.EntriesMatched)
package metrics
