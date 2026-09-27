package config

import (
	"os"
	"time"
)

// 第一版内置默认值。仅面向本地自用，生产环境应通过环境变量覆盖密钥。
const (
	defaultJWTSecret = "web2api-mvp-dev-secret"
	// EnvJWTSecret 是覆盖 JWT 密钥的环境变量名。
	EnvJWTSecret = "WEB2API_JWT_SECRET"
)

// Config 保存 Web2API 后端的基础配置。
type Config struct {
	// Port HTTP 服务监听端口。
	Port int
	// DataDir JSON 数据文件所在目录。
	DataDir string
	// JWTSecret 管理员 JWT 的签名密钥。
	JWTSecret string
	// JWTExpiry 管理员 JWT 有效期。
	JWTExpiry time.Duration
	// APIRateLimitRPM 每个 API Key 每分钟最大请求数（第一版固定值）。
	LogLevel        string
	IsProd          bool
	APIRateLimitRPM int
	// AccountWaitTimeout 无空闲账号时请求的最长等待时间。
	AccountWaitTimeout time.Duration
}

// Default 返回第一版 MVP 的默认配置。
// JWTSecret 优先取环境变量，便于不改代码切换密钥。
func Default() Config {
	return Config{
		Port:               8080,
		DataDir:            "data",
		IsProd:             false,
		LogLevel:           "debug",
		JWTSecret:          jwtSecret(),
		JWTExpiry:          24 * time.Hour,
		APIRateLimitRPM:    10,
		AccountWaitTimeout: 30 * time.Second,
	}
}

func jwtSecret() string {
	if v := os.Getenv(EnvJWTSecret); v != "" {
		return v
	}
	return defaultJWTSecret
}
