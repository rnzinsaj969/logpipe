// Package cap provides a per-service entry-count ceiling over a sliding time
// window. Entries that exceed the configured maximum for a service within the
// window are dropped, protecting downstream consumers from log bursts.
//
// Usage:
//
//	capper, err := cap.New(cap.Options{
//		Max:    100,
//		Window: time.Minute,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if capper.Allow(entry) {
//		// forward entry
//	}
package cap
