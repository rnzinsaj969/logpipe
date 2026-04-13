package config

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, cfg Config) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logpipe-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadValidConfig(t *testing.T) {
	cfg := Config{
		Sources: []Source{{Name: "app", Path: "/var/log/app.log"}},
		Level:   "info",
		Format:  "text",
	}
	path := writeTempConfig(t, cfg)

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Sources[0].Name != "app" {
		t.Errorf("expected source name %q, got %q", "app", loaded.Sources[0].Name)
	}
	if loaded.Format != "text" {
		t.Errorf("expected format %q, got %q", "text", loaded.Format)
	}
}

func TestLoadMissingSources(t *testing.T) {
	cfg := Config{Format: "json"}
	path := writeTempConfig(t, cfg)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing sources, got nil")
	}
}

func TestLoadInvalidFormat(t *testing.T) {
	cfg := Config{
		Sources: []Source{{Name: "svc", Path: "/tmp/svc.log"}},
		Format:  "yaml",
	}
	path := writeTempConfig(t, cfg)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestLoadInvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestLoadSourceMissingName(t *testing.T) {
	cfg := Config{
		Sources: []Source{{Path: "/tmp/svc.log"}},
	}
	path := writeTempConfig(t, cfg)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for source missing name, got nil")
	}
}
