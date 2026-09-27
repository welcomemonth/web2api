package storage

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

type testDoc struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestJSONStoreSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doc.json")
	s := NewJSONStore[testDoc](path)

	in := testDoc{ID: 1, Name: "a"}
	if err := s.Save(in); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != in {
		t.Errorf("roundtrip = %+v, want %+v", got, in)
	}
}

func TestJSONStoreLoadMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s := NewJSONStore[testDoc](path)

	if _, err := s.Load(); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("want os.ErrNotExist, got %v", err)
	}
}

func TestJSONStoreLoadCorrupt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doc.json")
	if err := os.WriteFile(path, []byte("{not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewJSONStore[testDoc](path)

	if _, err := s.Load(); err == nil {
		t.Errorf("want error for corrupt JSON, got nil")
	}
}

func TestJSONStoreSaveOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doc.json")
	s := NewJSONStore[testDoc](path)

	if err := s.Save(testDoc{ID: 1, Name: "first"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(testDoc{ID: 2, Name: "second"}); err != nil {
		t.Fatal(err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "second" {
		t.Errorf("got %q, want second（rename 覆盖应生效）", got.Name)
	}
}

func TestJSONStoreSaveCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deeper")
	s := NewJSONStore[testDoc](filepath.Join(dir, "doc.json"))

	if err := s.Save(testDoc{ID: 1}); err != nil {
		t.Fatalf("Save should create parent dirs: %v", err)
	}
}

func TestJSONStoreConcurrentSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doc.json")
	s := NewJSONStore[testDoc](path)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := s.Save(testDoc{ID: id}); err != nil {
				t.Errorf("Save(%d): %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load after concurrent saves: %v", err)
	}
	if got.ID < 0 || got.ID >= 50 {
		t.Errorf("load = %+v, want one of the written docs", got)
	}
}
