package rotate

import (
	"context"
	"os"
	"sync"
	"time"
)

// Watcher monitors a file path and signals when the file has been rotated
// (i.e. replaced by a new inode or truncated).
type Watcher struct {
	path     string
	interval time.Duration
	mu       sync.Mutex
	lastIno  uint64
	lastSize int64
}

// New creates a Watcher for the given file path, polling at the specified interval.
func New(path string, interval time.Duration) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
	}
}

// Watch blocks until the file is rotated or the context is cancelled.
// It sends on the returned channel when rotation is detected.
func (w *Watcher) Watch(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{}, 1)

	if err := w.snapshot(); err != nil {
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.rotated() {
					chtreturn
		
}

func (w *Watcher) snapshot() error {
	info, err := os.Stat(w.path)
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastIno = inode(info)
	w.lastSize = info.Size()
	return nil
}

func (w *Watcher) rotated() bool {
	info, err := os.Stat(w.path)
	if err != nil {
		return true
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	currentIno := inode(info)
	if currentIno != w.lastIno {
		return true
	}
	if info.Size() < w.lastSize {
		return true
	}
	w.lastSize = info.Size()
	return false
}
