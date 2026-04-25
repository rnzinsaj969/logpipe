// Package watermark provides a concurrency-safe high-water mark tracker for
// log entry timestamps.
//
// A Watermark records the latest timestamp observed across a stream of
// [reader.LogEntry] values. It can be used to detect out-of-order or late
// entries and to coordinate time-based processing across multiple sources.
//
// Usage:
//
//	w := watermark.New()
//	_ = w.Advance(entry)        // update high-water mark
//	if w.Behind(laterEntry) {   // check whether an entry is late
//	    // handle late arrival
//	}
//	w.Reset()                   // clear the mark
package watermark
