// Package limiter implements a per-service concurrency limiter for log entry
// processing pipelines.
//
// A Limiter enforces a maximum number of simultaneously in-flight entries for
// each unique service name. Callers must pair every successful Acquire with a
// corresponding Release to avoid leaking slots.
//
// Example usage:
//
//	l, err := limiter.New(5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if l.Acquire(entry) {
//		defer l.Release(entry)
//		// process entry
//	}
package limiter
