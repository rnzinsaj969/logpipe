// Package checkpoint provides a lightweight, file-backed store for persisting
// per-source read positions (byte offset and inode) across logpipe restarts.
//
// A Store is safe for concurrent use. Call Flush to atomically write the
// current in-memory state to disk using a write-then-rename strategy so that
// the checkpoint file is never left in a partially written state.
//
// Typical usage:
//
//	cp, err := checkpoint.New("/var/lib/logpipe/checkpoint.json")
//	if err != nil { ... }
//	st := cp.Get("my-service")
//	// ... advance reader to st.Offset ...
//	cp.Set("my-service", checkpoint.State{Offset: newOffset, Inode: inode})
//	_ = cp.Flush()
package checkpoint
