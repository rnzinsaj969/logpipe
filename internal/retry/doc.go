// Package retry provides a simple exponential back-off retry mechanism for
// use within logpipe when performing operations that may transiently fail,
// such as opening a log file that has not yet been created or writing to an
// output sink that is temporarily unavailable.
//
// Usage:
//
//	r := retry.New(retry.Config{
//		MaxAttempts: 5,
//		BaseDelay:   200 * time.Millisecond,
//		MaxDelay:    10 * time.Second,
//	})
//
//	err := r.Do(ctx, func() error {
//		return doSomethingFallible()
//	})
package retry
