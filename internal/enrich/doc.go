// Package enrich provides log entry enrichment by attaching additional
// metadata fields to structured log entries at processing time.
//
// An Enricher is constructed via New and accepts a variadic list of
// Option values that define enrichment rules. Rules are applied in
// order during each call to Apply.
//
// Supported enrichment options:
//
//   - WithStaticField: injects a constant key/value pair into the
//     entry's Fields map.
//   - WithServiceOverride: replaces the Service field of every entry
//     with the supplied value, useful when aggregating logs from a
//     sidecar or proxy that may report incorrect service names.
//
// Apply never mutates the original entry; it returns a shallow copy
// with a new Fields map so callers can safely pass entries through
// multiple pipeline stages concurrently.
package enrich
