// Package scope provides namespace-scoped access to nested fields
// within a LogEntry's Extra map.
//
// A Scoper targets a single namespace key and exposes Extract and Embed
// operations so that processors can read and write a sub-map without
// touching the rest of the entry.
//
// Example usage:
//
//	s, err := scope.New("kubernetes")
//	if err != nil { ... }
//	fields := s.Extract(entry)       // read nested map
//	fields["pod"] = "web-abc123"
//	entry = s.Embed(entry, fields)   // write back
package scope
