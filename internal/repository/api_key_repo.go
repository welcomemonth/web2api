package repository

import (
	"errors"
	"os"
	"sync"

	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/storage"
)

var ErrAPIKeyNotFound = errors.New("API Key 不存在")

// APIKeyRepository 负责账号层的 CRUD 与持久化。
// 存储为 []model.APIKey 的 JSON 数组，路径由服务启动时传入。
// 内存缓存防止每次请求都读文件。
type APIKeyRepository struct {
	store *storage.JSONStore[[]model.APIKey]
	mu    sync.RWMutex
	keys  []model.APIKey
}

func NewAPIKeyRepository(path string) *APIKeyRepository {
	return &APIKeyRepository{store: storage.NewJSONStore[[]model.APIKey](path)}
}

func (r *APIKeyRepository) Load() error {
	keys, err := r.store.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			r.mu.Lock()
			r.keys = []model.APIKey{}
			r.mu.Unlock()
			return nil
		}
		return err
	}
	r.mu.Lock()
	r.keys = keys
	r.mu.Unlock()
	return nil
}

// All 返回全部 key 的副本。
func (r *APIKeyRepository) All() []model.APIKey {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]model.APIKey(nil), r.keys...)
}

// GetByID 按 ID 返回 key 副本；不存在返回 (nil, false)。
func (r *APIKeyRepository) GetByID(id string) (*model.APIKey, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.keys {
		if r.keys[i].ID == id {
			cp := r.keys[i]
			return &cp, true
		}
	}
	return nil, false
}

// Create 新增 key 并落盘，key.ID 必须唯一，重复抛 ErrAPIKeyNotFound。
func (r *APIKeyRepository) Create(k model.APIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ex := range r.keys {
		if ex.ID == k.ID {
			return ErrAPIKeyNotFound
		}
	}
	r.keys = append(r.keys, k)
	return r.store.Save(r.keys)
}

// Update 按 ID 覆盖，ID 不存在返回 ErrAPIKeyNotFound。
func (r *APIKeyRepository) Update(k model.APIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, ex := range r.keys {
		if ex.ID == k.ID {
			r.keys[i] = k
			return r.store.Save(r.keys)
		}
	}
	return ErrAPIKeyNotFound
}

// Delete 按 ID 删除，ID 不存在返回 ErrAPIKeyNotFound。
func (r *APIKeyRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, ex := range r.keys {
		if ex.ID == id {
			r.keys = append(r.keys[:i], r.keys[i+1:]...)
			return r.store.Save(r.keys)
		}
	}
	return ErrAPIKeyNotFound
}
