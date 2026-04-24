// Package reorder provides a time-window reordering buffer for log entries.
//
// Out-of-order entries that arrive within a configurable window are held and
// emitted in ascending timestamp order once the window is full or the maximum
// age of the oldest entry is exceeded.
//
// Basic usage:
//
//	r, err := reorder.New(reorder.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, e := range incoming {
//		if ready := r.Add(e); ready != nil {
//			process(ready)
//		}
//	}
//	// flush remaining entries at end of stream
//	process(r.Flush())
package reorder
