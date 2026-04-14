package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

// State holds the persisted read position for a single source file.
type State struct {
	Offset int64  `json:"offset"`
	Inode  uint64 `json:"inode"`
}

// Store persists and retrieves per-source read offsets so that logpipe can
// resume from where it left off after a restart.
type Store struct {
	mu   sync.Mutex
	path string
	data map[string]State
}

// New loads an existing checkpoint file at path, or starts with an empty
// store if the file does not yet exist.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]State),
	}
	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &s.data); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the stored State for the given source key.
// If no entry exists, a zero-value State is returned.
func (s *Store) Get(source string) State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[source]
}

// Set updates the State for the given source key in memory.
func (s *Store) Set(source string, st State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[source] = st
}

// Flush writes the current in-memory state to disk atomically.
func (s *Store) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
