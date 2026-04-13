package rotate

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatcherDetectsTruncation(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rotate-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, _ = f.WriteString("initial content\n")

	w := New(f.Name(), 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := w.Watch(ctx)

	time.Sleep(50 * time.Millisecond)
	// Truncate the file to simulate rotation
	_ = f.Truncate(0)

	select {
	case _, ok := <-ch:
		if !ok {
			t.Fatal("channel closed unexpectedly without signal")
		}
		// rotation detected
	case <-ctx.Done():
		t.Fatal("timeout: rotation not detected")
	}
}

func TestWatcherCancelsCleanly(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rotate-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := New(f.Name(), 20*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	ch := w.Watch(ctx)
	cancel()

	select {
	case <-ch:
		// channel closed after cancel — expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout: watcher did not stop after context cancel")
	}
}

func TestWatcherInvalidPath(t *testing.T) {
	w := New("/nonexistent/path/file.log", 20*time.Millisecond)
	ctx := context.Background()

	ch := w.Watch(ctx)

	select {
	case <-ch:
		// closed immediately because snapshot fails
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout: expected immediate close for invalid path")
	}
}
