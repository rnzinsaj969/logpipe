// Package rank provides a Ranker that scores log entries against a set of
// field-value rules. Each matching rule contributes its configured weight to
// a running total which is written into a named extra field on the entry.
//
// Typical usage:
//
//	r, err := rank.New("priority", []rank.Rule{
//		{Field: "level",   Value: "error", Score: 100},
//		{Field: "service", Value: "payments", Score: 20},
//	})
//	...
//	scored := r.Apply(entry)
package rank
