// Package suppress provides a Suppressor that prevents repeated identical log
// entries from flooding downstream consumers.
//
// Entries are keyed by (service, message) pair. Once an entry is seen, further
// occurrences are dropped until the configured cooldown duration has elapsed.
//
// Example usage:
//
//	s, err := suppress.New(suppress.DefaultOptions())
//	if err != nil { ... }
//	if s.Apply(entry) {
//	    // forward entry
//	}
//
// Call Evict periodically to release memory for keys that have aged out.
package suppress
