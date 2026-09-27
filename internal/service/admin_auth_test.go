package service

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/welcomemonth/web2api/internal/auth"
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/repository"
)

func newAuthService(t *testing.T, dataDir string) (*AdminAuthService, *auth.JWTService) {
	t.Helper()

	repo := repository.NewAdminRepository(filepath.Join(dataDir, "admin.json"))
	svc := NewAdminInitService(repo)
	if _, err := svc.Ensure(); err != nil {
		t.Fatalf("初始化管理员失败: %v", err)
	}

	jwtSvc := auth.NewJWTService("test-secret", time.Hour)
	return NewAdminAuthService(repo, jwtSvc), jwtSvc
}

func TestAdminAuthLoginSuccess(t *testing.T) {
	svc, jwtSvc := newAuthService(t, t.TempDir())

	token, err := svc.Login(DefaultAdminUsername, DefaultAdminPassword)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	claims, err := jwtSvc.Parse(token)
	if err != nil {
		t.Fatalf("签发的 Token 应可校验: %v", err)
	}
	if claims.Username != DefaultAdminUsername {
		t.Errorf("Username = %q, want %q", claims.Username, DefaultAdminUsername)
	}
}

func TestAdminAuthLoginWrongPassword(t *testing.T) {
	svc, _ := newAuthService(t, t.TempDir())

	if _, err := svc.Login(DefaultAdminUsername, "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("错误密码应返回 ErrInvalidCredentials, got %v", err)
	}
}

func TestAdminAuthLoginUnknownUser(t *testing.T) {
	svc, _ := newAuthService(t, t.TempDir())

	if _, err := svc.Login("nobody", DefaultAdminPassword); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("不存在的用户应返回 ErrInvalidCredentials, got %v", err)
	}
}

func TestAdminAuthLoginEmptyInput(t *testing.T) {
	svc, _ := newAuthService(t, t.TempDir())

	for _, c := range []struct{ u, p string }{
		{"", ""},
		{DefaultAdminUsername, ""},
		{"", DefaultAdminPassword},
	} {
		if _, err := svc.Login(c.u, c.p); !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("Login(%q,%q) 应返回 ErrInvalidCredentials, got %v", c.u, c.p, err)
		}
	}
}

func TestAdminAuthLoginWithoutAdmin(t *testing.T) {
	// 尚未初始化管理员（admin.json 不存在）
	repo := repository.NewAdminRepository(filepath.Join(t.TempDir(), "admin.json"))
	svc := NewAdminAuthService(repo, auth.NewJWTService("test-secret", time.Hour))

	if _, err := svc.Login("admin", "admin123"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("无管理员时应返回 ErrInvalidCredentials, got %v", err)
	}
}

func TestAdminAuthLoginReflectsPasswordChange(t *testing.T) {
	dir := t.TempDir()
	repo := repository.NewAdminRepository(filepath.Join(dir, "admin.json"))
	if err := repo.Save(&model.Admin{Username: "admin", Password: "newpass"}); err != nil {
		t.Fatal(err)
	}
	svc := NewAdminAuthService(repo, auth.NewJWTService("test-secret", time.Hour))

	if _, err := svc.Login("admin", "newpass"); err != nil {
		t.Errorf("新密码应可登录: %v", err)
	}
	if _, err := svc.Login("admin", DefaultAdminPassword); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("旧密码应失效, got %v", err)
	}
}
