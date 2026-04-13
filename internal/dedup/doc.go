// Package dedup provides log entry deduplication for logpipe.
//
// A Deduplicator suppresses repeated log lines that share the same service,
// level, and message within a configurable time window. This prevents noisy,
// high-frequency duplicate entries from overwhelming downstream consumers.
//
// Usage:
//
//	dd := dedup.New(5 * time.Second)
//	if !dd.IsDuplicate(dedup.Entry{Service: "api", Level: "error", Message: "timeout"}) {
//	    // forward entry
//	}
//
// Call Evict periodically (e.g. via a ticker) to reclaim memory used by
// entries that have aged out of the dedup window.
package dedup
