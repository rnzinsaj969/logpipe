// Package normalize provides field-level normalisation for log entries.
//
// A Normalizer applies a configurable set of transformations to each
// reader.LogEntry it processes:
//
//   - Lowercase the level field for consistent downstream filtering.
//   - Trim leading/trailing whitespace from message and service fields.
//   - Fill in default values for empty level or service fields.
//
// The Apply method is pure: it never mutates the original entry and
// always returns a new value.
//
// Example:
//
//	n := normalize.New(normalize.DefaultOptions())
//	clean := n.Apply(rawEntry)
package normalize
