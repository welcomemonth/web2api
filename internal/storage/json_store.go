package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// JSONStore 是针对单个 JSON 文件的通用持久化存储。
// T 为落盘数据结构的类型。写操作通过互斥锁串行；读取依赖“临时文件 + rename”
// 的原子替换，读到的是完整的新旧文件之一，而非半写状态。
type JSONStore[T any] struct {
	path string
	mu   sync.Mutex
}

// NewJSONStore 创建针对 path 的 JSON 存储。
func NewJSONStore[T any](path string) *JSONStore[T] {
	return &JSONStore[T]{path: path}
}

// Load 读取并解析 JSON 文件。
// 文件不存在时返回错误，调用方可用 errors.Is(err, os.ErrNotExist) 判断并做初始化。
func (s *JSONStore[T]) Load() (T, error) {
	var out T

	data, err := os.ReadFile(s.path)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("解析 %s 失败: %w", s.path, err)
	}
	return out, nil
}

// Save 将 v 序列化为 JSON 并原子写入（写临时文件 → rename），保证并发写安全。
func (s *JSONStore[T]) Save(v T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 %s 失败: %w", s.path, err)
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*.json")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpName := tmp.Name()
	// 失败路径清理临时文件；rename 成功后 tmpName 已不存在，Remove 为无害 no-op。
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("同步临时文件失败: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}

	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("原子替换 %s 失败: %w", s.path, err)
	}
	return nil
}
