// Package suppress provides a Suppressor that prevents repeated identical log
// entries from flooding downstream consumers.
//
// Entries are keyed by (service, message) pair. Once an entry is seen, further
// occurrences are dropped until the configured cooldown duration has elapsed.
// The suppression count is tracked so callers can emit a summary such as
// "suppressed N identical messages" when the cooldown expires.
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
// A common pattern is to run Evict in a background goroutine:
//
//	go func() {
//		ticker := time.NewTicker(5 * time.Minute)
//		defer ticker.Stop()
//		for range ticker.C {
//			s.Evict()
//		}
//	}()
package suppress
