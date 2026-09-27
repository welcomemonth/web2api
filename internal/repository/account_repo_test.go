package repository

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/welcomemonth/web2api/internal/model"
)

func sampleAccount(id string) model.Account {
	return model.Account{
		ID:          id,
		Username:    "user-" + id,
		Password:    "pw",
		AccessToken: "tok",
		Status:      model.AccountStatusNormal,
		CreatedAt:   time.Now(),
	}
}

func TestAccountRepositoryLoadMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "accounts.json")
	repo := NewAccountRepository(path)

	if err := repo.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := repo.All(); len(got) != 0 {
		t.Errorf("文件不存在应初始化为空集合，got %d", len(got))
	}
}

func TestAccountRepositoryCRUD(t *testing.T) {
	path := filepath.Join(t.TempDir(), "accounts.json")
	repo := NewAccountRepository(path)
	if err := repo.Load(); err != nil {
		t.Fatal(err)
	}

	a1 := sampleAccount("1")
	a2 := sampleAccount("2")
	if err := repo.Create(a1); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(a2); err != nil {
		t.Fatal(err)
	}
	if got := repo.All(); len(got) != 2 {
		t.Fatalf("Create 后应有 2 个账号，got %d", len(got))
	}

	a1.Status = model.AccountStatusDisabled
	if err := repo.Update(a1); err != nil {
		t.Fatal(err)
	}
	got1, ok := repo.GetByID("1")
	if !ok || got1.Status != model.AccountStatusDisabled {
		t.Errorf("Update 后状态应为 disabled，got %+v", got1)
	}

	if err := repo.Delete("1"); err != nil {
		t.Fatal(err)
	}
	if got := repo.All(); len(got) != 1 || got[0].ID != "2" {
		t.Errorf("Delete 后应只剩账号 2，got %+v", got)
	}
}

func TestAccountRepositoryPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "accounts.json")
	repo := NewAccountRepository(path)
	if err := repo.Load(); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(sampleAccount("1")); err != nil {
		t.Fatal(err)
	}

	// 新仓库从磁盘加载，应恢复账号（模拟重启）
	repo2 := NewAccountRepository(path)
	if err := repo2.Load(); err != nil {
		t.Fatal(err)
	}
	if got := repo2.All(); len(got) != 1 || got[0].ID != "1" {
		t.Errorf("重启后应恢复账号，got %+v", got)
	}
}

func TestAccountRepositoryNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "accounts.json")
	repo := NewAccountRepository(path)
	if err := repo.Load(); err != nil {
		t.Fatal(err)
	}

	if err := repo.Update(sampleAccount("nope")); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Update 不存在应返回 ErrAccountNotFound，got %v", err)
	}
	if err := repo.Delete("nope"); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Delete 不存在应返回 ErrAccountNotFound，got %v", err)
	}
}

func TestAccountRepositoryAllCopyIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "accounts.json")
	repo := NewAccountRepository(path)
	if err := repo.Load(); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(sampleAccount("1")); err != nil {
		t.Fatal(err)
	}

	all := repo.All()
	all[0].Username = "mutated"

	if got := repo.All(); got[0].Username != "user-1" {
		t.Errorf("修改 All 返回值影响了内部缓存: %+v", got[0])
	}
}
