// Package expire provides a log entry processor that discards entries whose
// timestamp falls outside a configurable maximum age window.
//
// Usage:
//
//	p, err := expire.New(expire.Options{MaxAge: 5 * time.Minute})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if p.Apply(entry) {
//		// entry is recent enough to keep
//	}
//
// A zero timestamp is always retained so that entries without timing
// information are never silently discarded.
package expire
