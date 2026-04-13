package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempMain(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "logpipe.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestRunInvalidConfig(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()

	os.Args = []string{"logpipe", "/nonexistent/path/logpipe.yaml"}

	if err := run(); err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRunMissingSources(t *testing.T) {
	path := writeTempMain(t, "output:\n  format: text\n")

	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"logpipe", path}

	if err := run(); err == nil {
		t.Fatal("expected error for config with no sources, got nil")
	}
}
