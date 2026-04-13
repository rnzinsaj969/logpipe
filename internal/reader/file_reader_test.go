package reader_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/reader"
)

func TestFileReaderParsesEntry(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "fr-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fr, err := reader.NewFileReader(f.Name(), "mysvc")
	if err != nil {
		t.Fatal(err)
	}
	defer fr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entries := fr.Next(ctx)

	_, _ = f.WriteString("{\"level\":\"warn\",\"msg\":\"disk full\"}\n")

	select {
	case entry := <-entries:
		if entry.Service != "mysvc" {
			t.Errorf("expected service 'mysvc', got %q", entry.Service)
		}
		if entry.Message != "disk full" {
			t.Errorf("expected message 'disk full', got %q", entry.Message)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for entry")
	}
}

func TestFileReaderInvalidPath(t *testing.T) {
	_, err := reader.NewFileReader("/no/such/file.log", "svc")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
