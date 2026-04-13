// Package config provides configuration loading and validation for logpipe.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Source represents a named log source with a file path.
type Source struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Config holds the top-level logpipe configuration.
type Config struct {
	Sources  []Source `json:"sources"`
	Level    string   `json:"level"`
	Format   string   `json:"format"`
	Output   string   `json:"output"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required fields are present and values are acceptable.
func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("config: at least one source is required")
	}
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("config: source[%d] missing name", i)
		}
		if s.Path == "" {
			return fmt.Errorf("config: source[%d] missing path", i)
		}
	}
	validFormats := map[string]bool{"text": true, "json": true, "": true}
	if !validFormats[c.Format] {
		return fmt.Errorf("config: unsupported format %q (must be \"text\" or \"json\")", c.Format)
	}
	return nil
}
