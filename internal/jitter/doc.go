// Package jitter provides a Jitterer that applies a small, bounded random
// offset to the timestamp of each log entry.
//
// This is useful when multiple log sources emit entries with identical
// timestamps (e.g. millisecond resolution) and a downstream stage requires
// strict ordering. By spreading timestamps within a configurable window the
// risk of tie-breaking ambiguity is reduced without altering the relative
// ordering of entries that already carry distinct timestamps.
//
// Usage:
//
//	j, err := jitter.New(5*time.Millisecond, nil)
//	if err != nil { /* handle */ }
//	spread := j.Apply(entry)
package jitter
