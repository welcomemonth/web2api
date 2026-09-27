// Package auth 提供管理员认证所需的 JWT 签发与校验。
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken 表示 Token 缺失、签名不符、算法不符或已过期。
// 对外统一使用该错误，避免暴露具体失败原因。
var ErrInvalidToken = errors.New("无效的 Token")

// Claims 是管理员 JWT 的载荷。第一版仅携带用户名，不含角色与权限。
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTService 负责签发与校验管理员 Token。
type JWTService struct {
	secret []byte
	expiry time.Duration
}

// NewJWTService 创建 JWT 服务。expiry 为 Token 有效期（第一版 24 小时）。
func NewJWTService(secret string, expiry time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), expiry: expiry}
}

// Generate 为指定管理员签发 Token。
func (s *JWTService) Generate(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("签发 Token 失败: %w", err)
	}
	return signed, nil
}

// Parse 校验 Token 并返回其载荷。
// 任何失败（签名不符、算法被替换、过期等）都统一返回 ErrInvalidToken。
func (s *JWTService) Parse(token string) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		// 只接受 HMAC 签名，防止 alg=none 或非对称算法替换攻击。
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名算法: %v", t.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
