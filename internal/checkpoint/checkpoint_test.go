package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNewCreatesEmptyStore(t *testing.T) {
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Get("svc")
	if got.Offset != 0 || got.Inode != 0 {
		t.Fatalf("expected zero state, got %+v", got)
	}
}

func TestSetAndGet(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	s.Set("svc", checkpoint.State{Offset: 42, Inode: 7})
	got := s.Get("svc")
	if got.Offset != 42 || got.Inode != 7 {
		t.Fatalf("expected {42 7}, got %+v", got)
	}
}

func TestFlushAndReload(t *testing.T) {
	path := tempPath(t)
	s, _ := checkpoint.New(path)
	s.Set("alpha", checkpoint.State{Offset: 100, Inode: 3})
	s.Set("beta", checkpoint.State{Offset: 200, Inode: 5})
	if err := s.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	s2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := s2.Get("alpha"); got.Offset != 100 || got.Inode != 3 {
		t.Fatalf("alpha mismatch: %+v", got)
	}
	if got := s2.Get("beta"); got.Offset != 200 || got.Inode != 5 {
		t.Fatalf("beta mismatch: %+v", got)
	}
}

func TestNewInvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := checkpoint.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFlushIsAtomic(t *testing.T) {
	path := tempPath(t)
	s, _ := checkpoint.New(path)
	s.Set("x", checkpoint.State{Offset: 1})
	if err := s.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp file should not remain after flush")
	}
}
