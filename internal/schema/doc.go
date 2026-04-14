// Package schema defines field validation and normalization rules for
// structured log entries in logpipe.
//
// Use New to create a Validator with a custom set of required fields, or
// reference the package-level Required slice for the default schema.
// Normalize applies default values (e.g. timestamp, lowercase level) to an
// entry without mutating the original map.
package schema
