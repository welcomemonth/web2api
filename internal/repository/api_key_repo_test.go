package repository

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/welcomemonth/web2api/internal/model"
)

func defaultKey() model.APIKey {
	return model.APIKey{ID: "sk-1", Name: "key1", Key: "sk-1234", Status: model.APIKeyStatusEnabled, CreatedAt: time.Now()}
}

func TestAPIKeyRepoLoadMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_keys.json")
	repo := NewAPIKeyRepository(path)
	if err := repo.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := repo.All(); len(got) != 0 {
		t.Errorf("空文件应返回空 slice，got %d", len(got))
	}
}

func TestAPIKeyCRUD(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_keys.json")
	repo := NewAPIKeyRepository(path)
	repo.Load()

	// create
	if err := repo.Create(defaultKey()); err != nil {
		t.Fatal(err)
	}
	if got := repo.All(); len(got) != 1 {
		t.Fatalf("create 1, got %d", len(got))
	}

	// update
	k := defaultKey()
	k.Status = model.APIKeyStatusDisabled
	if err := repo.Update(k); err != nil {
		t.Fatal(err)
	}
	if got, _ := repo.GetByID("sk-1"); got.Status != model.APIKeyStatusDisabled {
		t.Errorf("update status expected disabled, got %+v", got)
	}

	// delete
	if err := repo.Delete("sk-1"); err != nil {
		t.Fatal(err)
	}
	if got := repo.All(); len(got) != 0 {
		t.Errorf("delete should leave empty, got %+v", got)
	}
}

func TestAPIKeyPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_keys.json")
	repo := NewAPIKeyRepository(path)
	repo.Load()
	repo.Create(defaultKey())

	// reload from disk
	repo2 := NewAPIKeyRepository(path)
	repo2.Load()
	if got := repo2.All(); len(got) != 1 || got[0].ID != "sk-1" {
		t.Errorf("persist failed, got %+v", got)
	}
}

func TestAPIKeyNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_keys.json")
	repo := NewAPIKeyRepository(path)
	repo.Load()

	if err := repo.Update(defaultKey()); !errors.Is(err, ErrAPIKeyNotFound) {
		t.Errorf("update non-exist should ErrAPIKeyNotFound, got %v", err)
	}
	if err := repo.Delete("nope"); !errors.Is(err, ErrAPIKeyNotFound) {
		t.Errorf("delete non-exist should ErrAPIKeyNotFound, got %v", err)
	}
}

func TestAPIKeySliceCopyIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api_keys.json")
	repo := NewAPIKeyRepository(path)
	repo.Create(defaultKey())

	all := repo.All()
	all[0].Name = "mutated"

	if got := repo.All(); got[0].Name != "key1" {
		t.Errorf("修改 All 返回值影响内部缓存: %+v", got[0])
	}
}
