// Package dedupe implements field-level deduplication for structured log entries.
//
// A Deduplicator tracks recently seen log entries by hashing a configurable
// subset of their fields (message, level, service). Entries whose key fields
// match a previously seen entry within the configured time window are flagged
// as duplicates and can be suppressed by the caller.
//
// Usage:
//
//	d := dedupe.New(dedupe.DefaultOptions())
//	if !d.IsDuplicate(entry) {
//	    output.Write(entry)
//	}
//
// Call Evict periodically to release memory held by stale entries.
package dedupe
