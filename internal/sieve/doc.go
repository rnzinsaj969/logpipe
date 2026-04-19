// Package sieve provides allow-rule based filtering for log entries.
//
// A Sieve holds one or more Rules. An entry passes the sieve if it satisfies
// ANY rule. When no rules are configured every entry is allowed through,
// making the sieve a transparent no-op.
//
// Each Rule can match on Level, Service, and a regular-expression Pattern
// applied to the message. All non-empty fields within a single rule must
// match simultaneously (AND semantics within a rule).
//
// Example:
//
//	s, err := sieve.New([]sieve.Rule{
//		{Level: "error"},
//		{Service: "payments", Pattern: "timeout"},
//	})
//	if err != nil { ... }
//	if s.Apply(entry) { /* keep */ }
package sieve
