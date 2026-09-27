package app

import (
	"fmt"
	"path/filepath"

	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/repository"
)

// App 是应用启动时装配好的依赖容器。
// 各 Repository 在构造时把 JSON 数据加载进内存，业务层运行期间优先访问内存数据。
type App struct {
	Config   *config.Config
	Admin    *repository.AdminRepository
	Accounts *repository.AccountRepository
	APIKeys  *repository.APIKeyRepository
}

// New 依据配置定位 DataDir 下的各 JSON 文件并加载到内存。
// 任一文件损坏都会返回错误，避免带病启动。
func New(cfg *config.Config) (*App, error) {
	a := &App{
		Config:   cfg,
		Admin:    repository.NewAdminRepository(filepath.Join(cfg.DataDir, "admin.json")),
		Accounts: repository.NewAccountRepository(filepath.Join(cfg.DataDir, "accounts.json")),
		APIKeys:  repository.NewAPIKeyRepository(filepath.Join(cfg.DataDir, "api_keys.json")),
	}
	if err := a.load(); err != nil {
		return nil, err
	}
	return a, nil
}

// load 按文件逐个加载；文件不存在时各 Repository 会初始化为空集合。
func (a *App) load() error {
	if err := a.Admin.Load(); err != nil {
		return fmt.Errorf("加载 admin.json 失败: %w", err)
	}
	if err := a.Accounts.Load(); err != nil {
		return fmt.Errorf("加载 accounts.json 失败: %w", err)
	}
	if err := a.APIKeys.Load(); err != nil {
		return fmt.Errorf("加载 api_keys.json 失败: %w", err)
	}
	return nil
}
