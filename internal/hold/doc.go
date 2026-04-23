// Package hold implements a conditional log-entry buffer.
//
// A Holder accumulates [reader.LogEntry] values until a caller-supplied
// Predicate returns true, at which point the entire buffer is flushed and
// returned to the caller. This is useful for implementing patterns such as
// "capture the last N debug lines and emit them only when an error occurs".
//
// # Basic usage
//
//	h, err := hold.New(100, func(e reader.LogEntry) bool {
//		return e.Level == "error"
//	})
//	if err != nil { ... }
//
//	for _, entry := range stream {
//		if entries, ok := h.Add(entry); ok {
//			for _, flushed := range entries {
//				output.Write(flushed)
//			}
//		}
//	}
package hold
