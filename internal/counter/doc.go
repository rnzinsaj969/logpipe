// Package counter provides a thread-safe per-service log entry counter
// for use in pipelines and aggregation layers.
//
// Usage:
//
//	c := counter.New()
//	c.Inc(entry.Service)
//	snap := c.Snapshot()
package counter
