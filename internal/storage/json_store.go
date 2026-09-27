package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type JSONStore struct {
	path        string
	defaultData any
	mu          sync.RWMutex
}

func NewJSONStore(path string, defaultData any) *JSONStore {
	return &JSONStore{path: path, defaultData: defaultData}
}

func (s *JSONStore) Path() string { return s.path }

func (s *JSONStore) Ensure() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeJSONFileLocked(s.path, s.defaultData)
}

func (s *JSONStore) LoadInto(v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		if err := writeJSONFileLocked(s.path, s.defaultData); err != nil {
			return err
		}
	}
	raw, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		raw, _ = json.Marshal(s.defaultData)
	}
	return json.Unmarshal(raw, v)
}

func (s *JSONStore) LoadAny() (any, error) {
	var v any
	if err := s.LoadInto(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *JSONStore) Save(v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return writeJSONFileLocked(s.path, v)
}

func writeJSONFileLocked(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
