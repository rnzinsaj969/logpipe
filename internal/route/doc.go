// Package route provides rule-based routing of log entries to named
// destinations.
//
// A Router holds an ordered list of Rules. Each Rule may constrain matching
// by log level, service name (via a regular expression), or both. The first
// rule whose predicates all match an entry determines the destination label
// returned by Match. A rule with no predicates acts as a catch-all.
//
// Example:
//
//	rt, err := route.New([]route.Rule{
//		{Destination: "critical", Level: "error"},
//		{Destination: "auth",     ServicePattern: "^auth"},
//		{Destination: "default"},
//	})
//	dst := rt.Match(entry) // "critical", "auth", or "default"
package route
