// Package counter provides a thread-safe per-service log entry counter
// for use in pipelines and aggregation layers.
//
// Usage:
//
//	c := counter.New()
//	c.Inc(entry.Service)
//	snap := c.Snapshot()
//
// Snapshot returns a copy of the current counts keyed by service name,
// safe to read without holding any lock. The original counter continues
// to accumulate entries independently of any snapshot taken.
//
// All methods on Counter are safe for concurrent use by multiple goroutines.
package counter
