package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/welcomemonth/web2api/internal/model"
)

func defaultAdmin() *model.Admin {
	return &model.Admin{Username: "admin", Password: "admin123"}
}

func TestAdminRepositoryEnsureAdminCreatesDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "admin.json")
	repo := NewAdminRepository(path)

	got, err := repo.EnsureAdmin(defaultAdmin())
	if err != nil {
		t.Fatalf("EnsureAdmin: %v", err)
	}
	if got == nil || got.Username != "admin" || got.Password != "admin123" {
		t.Errorf("got %+v, want default admin", got)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("admin.json 应已落盘: %v", err)
	}
}

func TestAdminRepositoryEnsureAdminLoadsExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "admin.json")

	// 先落盘一个已有管理员
	existing := &model.Admin{Username: "root", Password: "secret"}
	if err := NewAdminRepository(path).Save(existing); err != nil {
		t.Fatal(err)
	}

	// 新仓库 EnsureAdmin 应读取已有数据，而非用默认值覆盖
	repo := NewAdminRepository(path)
	got, err := repo.EnsureAdmin(defaultAdmin())
	if err != nil {
		t.Fatalf("EnsureAdmin: %v", err)
	}
	if got.Username != "root" || got.Password != "secret" {
		t.Errorf("got %+v, want existing admin（不应被默认值覆盖）", got)
	}
}

func TestAdminRepositoryEnsureAdminCorrupt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "admin.json")
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := NewAdminRepository(path)
	if _, err := repo.EnsureAdmin(defaultAdmin()); err == nil {
		t.Errorf("损坏的 admin.json 应返回错误，而非静默覆盖")
	}
}

func TestAdminRepositoryGetCopyIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "admin.json")
	repo := NewAdminRepository(path)
	if _, err := repo.EnsureAdmin(defaultAdmin()); err != nil {
		t.Fatal(err)
	}

	got := repo.Get()
	got.Password = "mutated" // 修改返回值副本

	again := repo.Get()
	if again.Password != "admin123" {
		t.Errorf("修改 Get 返回值影响了内部缓存: %+v", again)
	}
}
