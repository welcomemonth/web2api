package service

import (
	"crypto/subtle"
	"errors"

	"github.com/welcomemonth/web2api/internal/auth"
	"github.com/welcomemonth/web2api/internal/repository"
)

// ErrInvalidCredentials 表示用户名或密码不正确。
// 用户名不存在与密码错误返回同一个错误，避免暴露账号是否存在。
var ErrInvalidCredentials = errors.New("用户名或密码错误")

// AdminAuthService 负责管理员登录校验并签发 JWT。
type AdminAuthService struct {
	repo *repository.AdminRepository
	jwt  *auth.JWTService
}

// NewAdminAuthService 创建管理员认证服务。
func NewAdminAuthService(repo *repository.AdminRepository, jwtSvc *auth.JWTService) *AdminAuthService {
	return &AdminAuthService{repo: repo, jwt: jwtSvc}
}

// Login 校验用户名密码，成功则返回签发的 JWT。
// 校验失败统一返回 ErrInvalidCredentials。
func (s *AdminAuthService) Login(username, password string) (string, error) {
	admin := s.repo.Get()
	if admin == nil {
		return "", ErrInvalidCredentials
	}
	if !secureEqual(admin.Username, username) || !secureEqual(admin.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(admin.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// secureEqual 以恒定时间比较字符串，避免通过响应耗时逐字符猜测凭证。
func secureEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
