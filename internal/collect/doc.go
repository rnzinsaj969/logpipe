// Package collect provides a thread-safe, bounded log entry accumulator.
//
// A Collector buffers up to a configurable maximum number of entries.
// When the limit is reached the oldest entry is evicted to make room for
// the newest, giving ring-buffer semantics.
//
// Typical usage:
//
//	c, err := collect.New(500)
//	if err != nil { ... }
//	c.Add(entry)
//	entries := c.Flush() // snapshot and reset
package collect
