package service

import (
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/repository"
)

// 第一版内置的默认管理员。仅面向本地自用，密码明文存储。
const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
)

// AdminInitService 负责首次启动时确保存在一个可用管理员。
type AdminInitService struct {
	repo *repository.AdminRepository
}

// NewAdminInitService 创建默认管理员初始化服务。
func NewAdminInitService(repo *repository.AdminRepository) *AdminInitService {
	return &AdminInitService{repo: repo}
}

// Ensure 确保管理员已就绪：
// admin.json 不存在时创建默认管理员，已存在则直接复用，不覆盖已有密码。
func (s *AdminInitService) Ensure() (*model.Admin, error) {
	return s.repo.EnsureAdmin(&model.Admin{
		Username: DefaultAdminUsername,
		Password: DefaultAdminPassword,
	})
}
