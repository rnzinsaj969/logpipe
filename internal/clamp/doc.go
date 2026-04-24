// Package clamp provides a rate-limiting component that enforces a maximum
// number of log entries per service within a sliding time window.
//
// A Clamp is created with a maximum count and a window duration. Each call to
// Allow checks whether the given service has exceeded its quota for the current
// window. When the window expires the counter resets automatically.
//
// Example usage:
//
//	c, err := clamp.New(100, time.Minute)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if c.Allow(entry) {
//		// forward entry downstream
//	}
package clamp
