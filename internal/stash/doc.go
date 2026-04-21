// Package stash provides a thread-safe, keyed bucket store for log entries.
//
// Entries are grouped under arbitrary string keys and can be retrieved in
// bulk via Flush. Each bucket enforces a configurable capacity; when the
// limit is reached the oldest entry is evicted to make room for the newest
// (FIFO ring behaviour).
//
// Typical usage is to accumulate related entries — e.g. all lines belonging
// to a single request trace — and emit them together once a terminal event
// is observed.
package stash
