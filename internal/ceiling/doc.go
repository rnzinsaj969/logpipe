// Package ceiling provides a per-service entry rate limiter that enforces
// a maximum number of log entries allowed within a sliding time window.
//
// A Ceiling tracks counts independently for each service name. Once a service
// exceeds the configured maximum, subsequent entries are dropped until the
// window resets.
//
// Example usage:
//
//	c, err := ceiling.New(100, time.Minute)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if c.Allow(entry) {
//		// process entry
//	}
package ceiling
