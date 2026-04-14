// Package throttle implements per-service log entry throttling for logpipe.
//
// A Throttler limits the number of log entries accepted from a named service
// within a configurable sliding time window. Entries that exceed the limit are
// signalled for dropping by returning false from Allow.
//
// Usage:
//
//	th := throttle.New(time.Second, 100, nil)
//	if th.Allow(entry.Service) {
//		// forward entry downstream
//	}
//
// The zero window or zero maxCount is valid; a zero maxCount blocks all
// entries. The Clock dependency is injectable for deterministic testing.
package throttle
