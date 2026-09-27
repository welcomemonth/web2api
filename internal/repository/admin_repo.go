package repository

import (
	"errors"
	"os"
	"sync"

	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/storage"
)

// AdminRepository 管理单个管理员账号的读取、保存与首次初始化。
// 内存中缓存当前管理员，业务层优先访问内存数据。
type AdminRepository struct {
	store *storage.JSONStore[model.Admin]
	mu    sync.RWMutex
	admin *model.Admin
}

// NewAdminRepository 创建针对 path（admin.json）的管理员仓库。
func NewAdminRepository(path string) *AdminRepository {
	return &AdminRepository{store: storage.NewJSONStore[model.Admin](path)}
}

// Get 返回内存中缓存的管理员副本；尚未初始化时返回 nil。
func (r *AdminRepository) Get() *model.Admin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.admin == nil {
		return nil
	}
	cp := *r.admin
	return &cp
}

// Save 持久化管理员并更新内存缓存。
func (r *AdminRepository) Save(a *model.Admin) error {
	if a == nil {
		return errors.New("admin 不能为空")
	}
	if err := r.store.Save(*a); err != nil {
		return err
	}
	cp := *a
	r.mu.Lock()
	r.admin = &cp
	r.mu.Unlock()
	return nil
}

// EnsureAdmin 首次启动时初始化管理员：
// admin.json 已存在则读取并缓存；不存在则用 defaultAdmin 创建并落盘。
// 损坏等其它错误原样返回，不做静默覆盖。
func (r *AdminRepository) EnsureAdmin(defaultAdmin *model.Admin) (*model.Admin, error) {
	a, err := r.store.Load()
	if err == nil {
		r.mu.Lock()
		r.admin = &a
		r.mu.Unlock()
		return r.Get(), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err := r.Save(defaultAdmin); err != nil {
		return nil, err
	}
	return r.Get(), nil
}
