package service

import (
	"path/filepath"
	"testing"

	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/repository"
)

func newInitService(dataDir string) *AdminInitService {
	return NewAdminInitService(repository.NewAdminRepository(filepath.Join(dataDir, "admin.json")))
}

func TestAdminInitServiceCreatesDefault(t *testing.T) {
	dir := t.TempDir()
	svc := newInitService(dir)

	got, err := svc.Ensure()
	if err != nil {
		t.Fatalf("Ensure: %v", err)
	}
	if got.Username != DefaultAdminUsername || got.Password != DefaultAdminPassword {
		t.Errorf("got %+v, want 默认管理员 admin/admin123", got)
	}
}

func TestAdminInitServiceIsIdempotent(t *testing.T) {
	dir := t.TempDir()

	if _, err := newInitService(dir).Ensure(); err != nil {
		t.Fatal(err)
	}
	// 第二次调用（模拟重启）不应报错，也不应改变账号
	got, err := newInitService(dir).Ensure()
	if err != nil {
		t.Fatalf("第二次 Ensure: %v", err)
	}
	if got.Username != DefaultAdminUsername {
		t.Errorf("got %+v, want 默认管理员", got)
	}
}

func TestAdminInitServiceKeepsExisting(t *testing.T) {
	dir := t.TempDir()

	// 预置一个改过密码的管理员
	repo := repository.NewAdminRepository(filepath.Join(dir, "admin.json"))
	if err := repo.Save(&model.Admin{Username: "root", Password: "changed"}); err != nil {
		t.Fatal(err)
	}

	got, err := newInitService(dir).Ensure()
	if err != nil {
		t.Fatalf("Ensure: %v", err)
	}
	if got.Username != "root" || got.Password != "changed" {
		t.Errorf("got %+v, want 保留已有管理员（不覆盖密码）", got)
	}
}
