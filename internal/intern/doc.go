// Package intern implements a concurrency-safe string interning pool.
//
// High-throughput log pipelines repeatedly encounter the same short strings
// for fields such as "level" ("info", "warn", "error") and "service" names.
// Interning ensures that only a single copy of each unique string is kept in
// memory, reducing GC pressure and improving cache locality.
//
// Basic usage:
//
//	pool := intern.New()
//	level := pool.Intern(entry.Level)   // returns canonical copy
//	service := pool.Intern(entry.Service)
//
// The pool is safe for concurrent use. Call Reset to release all interned
// strings when the pool is no longer needed.
package intern
