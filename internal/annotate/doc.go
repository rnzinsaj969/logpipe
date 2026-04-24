// Package annotate provides a rule-based annotation processor that adds
// structured key-value fields to log entries whose messages match a set of
// regular-expression patterns.
//
// Each Rule pairs a compiled pattern with a target Extra field key and value.
// When the pattern matches the entry's Message, the key-value pair is merged
// into the entry's Extra map.  Multiple rules are evaluated in order and all
// matching rules are applied.  The original entry is never mutated.
//
// Example:
//
//	a, err := annotate.New([]annotate.Rule{
//		{Pattern: `timeout`, Key: "category", Value: "network"},
//		{Pattern: `panic`,   Key: "severity", Value: "critical"},
//	})
//	if err != nil { /* handle */ }
//	annotated := a.Apply(entry)
package annotate
