// Package label provides rule-based metadata labelling for log entries.
//
// A Labeler holds a slice of Rule values. Each Rule specifies a field to
// inspect, a substring to search for, and a label key/value pair to attach
// when the rule matches. Rules are evaluated in order and all matching rules
// are applied — the entry is never mutated; Apply always returns a copy.
//
// Example usage:
//
//	l := label.New([]label.Rule{
//		{Field: "level",   Contains: "error", Label: "severity", Value: "high"},
//		{Field: "message", Contains: "timeout", Label: "alert",   Value: "true"},
//	})
//	labelled := l.Apply(entry)
package label
