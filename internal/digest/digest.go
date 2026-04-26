// Package digest computes a short fingerprint for a log entry based on
// selected fields. The fingerprint is suitable for deduplication keys,
// cache lookups, and audit correlation.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/logpipe/logpipe/internal/reader"
)

// Options controls which fields contribute to the digest.
type Options struct {
	// IncludeMessage includes the entry message in the hash.
	IncludeMessage bool
	// IncludeLevel includes the log level in the hash.
	IncludeLevel bool
	// IncludeService includes the service name in the hash.
	IncludeService bool
	// ExtraKeys lists keys from the Extra map to include.
	ExtraKeys []string
}

// DefaultOptions returns a sensible default that hashes message, level, and
// service.
func DefaultOptions() Options {
	return Options{
		IncludeMessage: true,
		IncludeLevel:   true,
		IncludeService: true,
	}
}

// Digester computes fingerprints for log entries.
type Digester struct {
	opts Options
}

// New returns a Digester configured with opts.
func New(opts Options) (*Digester, error) {
	if !opts.IncludeMessage && !opts.IncludeLevel && !opts.IncludeService && len(opts.ExtraKeys) == 0 {
		return nil, fmt.Errorf("digest: at least one field must be included")
	}
	return &Digester{opts: opts}, nil
}

// Sum returns a hex-encoded SHA-256 digest of the configured fields of e.
func (d *Digester) Sum(e reader.LogEntry) string {
	h := sha256.New()
	if d.opts.IncludeMessage {
		_, _ = fmt.Fprint(h, "msg:", e.Message, "\x00")
	}
	if d.opts.IncludeLevel {
		_, _ = fmt.Fprint(h, "lvl:", e.Level, "\x00")
	}
	if d.opts.IncludeService {
		_, _ = fmt.Fprint(h, "svc:", e.Service, "\x00")
	}
	keys := make([]string, len(d.opts.ExtraKeys))
	copy(keys, d.opts.ExtraKeys)
	sort.Strings(keys)
	for _, k := range keys {
		v, ok := e.Extra[k]
		if !ok {
			continue
		}
		_, _ = fmt.Fprintf(h, "extra:%s=%v\x00", k, v)
	}
	return hex.EncodeToString(h.Sum(nil))
}
