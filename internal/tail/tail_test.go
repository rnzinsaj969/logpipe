package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/tail"
)

func TestTailerReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tailer, err := tail.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer tailer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lines := tailer.Lines(ctx)

	_, _ = f.WriteString("{\"level\":\"info\",\"msg\":\"hello\"}\n")

	select {
	case line := <-lines:
		if line == "" {
			t.Error("expected non-empty line")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for line")
	}
}

func TestTailerClosesOnContextCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cancel-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tailer, err := tail.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer tailer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	lines := tailer.Lines(ctx)
	cancel()

	// Channel should be closed shortly after cancellation.
	select {
	case _, ok := <-lines:
		if ok {
			t.Error("expected channel to be closed after context cancel")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}

func TestTailerInvalidPath(t *testing.T) {
	_, err := tail.New("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
