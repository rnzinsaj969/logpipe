// Package truncate caps the byte length of log entry fields to prevent
// oversized payloads from propagating through the logpipe pipeline.
//
// Usage:
//
//	tr := truncate.New(truncate.Options{
//		MaxMessageBytes: 4096,
//		MaxFieldBytes:   512,
//	})
//
//	if tr.Needed(entry) {
//		entry = tr.Apply(entry)
//	}
//
// Apply never mutates the original entry; it returns a shallow copy with
// string fields replaced where they exceed the configured limits. Truncated
// values are suffixed with the Unicode ellipsis character (…) so consumers
// can detect data loss.
package truncate
