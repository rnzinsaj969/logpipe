// Package batch provides a size- and age-based log entry batcher.
//
// A Batcher accumulates [reader.LogEntry] values and signals when the batch
// is ready to be flushed — either because it has reached a maximum number of
// entries (MaxSize) or because the oldest entry has been held for longer than
// a configured duration (MaxAge).
//
// Typical usage:
//
//	b := batch.New(batch.DefaultOptions())
//	for _, e := range entries {
//		if b.Add(e) {
//			flush(b.Flush())
//		}
//	}
//	// flush any remaining entries
//	if b.Len() > 0 {
//		flush(b.Flush())
//	}
package batch
