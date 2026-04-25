// Package ceiling provides a per-service rolling-window entry cap.
//
// A Ceiling tracks how many log entries each service has emitted within a
// configurable time window. Once a service reaches the configured maximum,
// further entries are dropped until the oldest observations fall outside the
// window and free up capacity again.
//
// Typical usage:
//
//	ceil, err := ceiling.New(100, time.Minute)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if ceil.Allow(entry) {
//		// forward entry downstream
//	}
package ceiling
