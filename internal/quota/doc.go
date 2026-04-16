// Package quota provides per-service entry rate limiting over a sliding
// time window. It is designed to cap the number of log entries a single
// service may emit within a configurable duration, preventing any one
// source from overwhelming the pipeline.
//
// Usage:
//
//	q, err := quota.New(quota.Options{MaxEntries: 500, Window: time.Minute})
//	if err != nil { ... }
//	if err := q.Allow(entry.Service); errors.Is(err, quota.ErrQuotaExceeded) {
//	    // drop or flag the entry
//	}
package quota
