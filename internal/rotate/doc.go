// Package rotate provides file rotation detection for log sources.
//
// A Watcher polls a file at a configurable interval and signals on a channel
// when the file appears to have been rotated — either because its inode has
// changed (the file was replaced) or because its size has decreased
// (the file was truncated in-place).
//
// Typical usage:
//
//	w := rotate.New("/var/log/app.log", 500*time.Millisecond)
//	ch := w.Watch(ctx)
//	select {
//	case <-ch:
//		// reopen the file
//	case <-ctx.Done():
//	}
package rotate
