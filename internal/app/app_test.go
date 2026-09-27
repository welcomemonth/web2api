package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/repository"
)

func testConfig(dataDir string) *config.Config {
	cfg := config.Default()
	cfg.DataDir = dataDir
	return &cfg
}

func TestNewLoadsExistingData(t *testing.T) {
	dir := t.TempDir()

	// 预置三份 JSON，模拟已有数据目录
	adminRepo := repository.NewAdminRepository(filepath.Join(dir, "admin.json"))
	if err := adminRepo.Save(&model.Admin{Username: "root", Password: "pw"}); err != nil {
		t.Fatal(err)
	}
	accountRepo := repository.NewAccountRepository(filepath.Join(dir, "accounts.json"))
	if err := accountRepo.Create(model.Account{
		ID: "1", Username: "u1", Status: model.AccountStatusNormal, CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}
	keyRepo := repository.NewAPIKeyRepository(filepath.Join(dir, "api_keys.json"))
	if err := keyRepo.Create(model.APIKey{
		ID: "k1", Name: "n1", Key: "sk-abc", Status: model.APIKeyStatusEnabled, CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	application, err := New(testConfig(dir))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if got := application.Admin.Get(); got == nil || got.Username != "root" {
		t.Errorf("admin 未恢复: %+v", got)
	}
	if got := application.Accounts.All(); len(got) != 1 || got[0].ID != "1" {
		t.Errorf("accounts 未恢复: %+v", got)
	}
	if got := application.APIKeys.All(); len(got) != 1 || got[0].ID != "k1" {
		t.Errorf("api_keys 未恢复: %+v", got)
	}
}

func TestNewMissingFilesIsEmpty(t *testing.T) {
	application, err := New(testConfig(filepath.Join(t.TempDir(), "empty")))
	if err != nil {
		t.Fatalf("文件缺失时 New 不应报错: %v", err)
	}

	if got := application.Admin.Get(); got != nil {
		t.Errorf("admin 应为未初始化，got %+v", got)
	}
	if got := application.Accounts.All(); len(got) != 0 {
		t.Errorf("accounts 应为空，got %+v", got)
	}
	if got := application.APIKeys.All(); len(got) != 0 {
		t.Errorf("api_keys 应为空，got %+v", got)
	}
}

func TestNewCorruptFileFails(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "accounts.json"), []byte("{bad"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := New(testConfig(dir)); err == nil {
		t.Errorf("损坏的 accounts.json 应导致启动失败")
	}
}
