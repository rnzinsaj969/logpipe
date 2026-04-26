// Package digest provides a Digester that computes a deterministic
// SHA-256 fingerprint for a log entry.
//
// The set of fields included in the hash is controlled by Options.
// At least one field must be enabled; constructing a Digester with an
// empty Options returns an error.
//
// Typical use:
//
//	d, err := digest.New(digest.DefaultOptions())
//	if err != nil { ... }
//	fingerprint := d.Sum(entry)
//
// The returned fingerprint is a lowercase hex-encoded SHA-256 string
// (64 characters). Timestamps and fields not listed in Options are
// intentionally excluded so that the same logical event always produces
// the same fingerprint.
package digest
