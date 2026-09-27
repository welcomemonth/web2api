package repository

import (
	"errors"
	"os"
	"sync"

	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/storage"
)

// ErrAccountNotFound 表示按 ID 未找到账号。
var ErrAccountNotFound = errors.New("账号不存在")

// AccountRepository 管理全部 Qwen 账号，内存持有账号集合并 CRUD 落盘 accounts.json。
type AccountRepository struct {
	store    *storage.JSONStore[[]model.Account]
	mu       sync.RWMutex
	accounts []model.Account
}

// NewAccountRepository 创建针对 path（accounts.json）的账号仓库。
func NewAccountRepository(path string) *AccountRepository {
	return &AccountRepository{store: storage.NewJSONStore[[]model.Account](path)}
}

// Load 从 accounts.json 加载全部账号到内存；文件不存在时初始化为空集合。
func (r *AccountRepository) Load() error {
	accs, err := r.store.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			r.mu.Lock()
			r.accounts = []model.Account{}
			r.mu.Unlock()
			return nil
		}
		return err
	}
	r.mu.Lock()
	r.accounts = accs
	r.mu.Unlock()
	return nil
}

// All 返回全部账号的副本。
func (r *AccountRepository) All() []model.Account {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]model.Account(nil), r.accounts...)
}

// GetByID 按 ID 返回账号副本；不存在时返回 (nil, false)。
func (r *AccountRepository) GetByID(id string) (*model.Account, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			cp := r.accounts[i]
			return &cp, true
		}
	}
	return nil, false
}

// Create 新增账号并落盘。
func (r *AccountRepository) Create(a model.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts = append(r.accounts, a)
	return r.store.Save(r.accounts)
}

// Update 按 ID 覆盖账号；ID 不存在时返回 ErrAccountNotFound。
func (r *AccountRepository) Update(a model.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		if r.accounts[i].ID == a.ID {
			r.accounts[i] = a
			return r.store.Save(r.accounts)
		}
	}
	return ErrAccountNotFound
}

// Delete 按 ID 删除账号；ID 不存在时返回 ErrAccountNotFound。
func (r *AccountRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			r.accounts = append(r.accounts[:i], r.accounts[i+1:]...)
			return r.store.Save(r.accounts)
		}
	}
	return ErrAccountNotFound
}
