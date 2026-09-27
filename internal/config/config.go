package config

import "time"

// Config 保存 Web2API 后端的基础配置。
type Config struct {
	// Port HTTP 服务监听端口。
	Port int
	// DataDir JSON 数据文件所在目录。
	DataDir string
	// JWTExpiry 管理员 JWT 有效期。
	JWTExpiry time.Duration
	// APIRateLimitRPM 每个 API Key 每分钟最大请求数（第一版固定值）。
	APIRateLimitRPM int
	// AccountWaitTimeout 无空闲账号时请求的最长等待时间。
	AccountWaitTimeout time.Duration
}

// Default 返回第一版 MVP 的默认配置。
func Default() Config {
	return Config{
		Port:               8080,
		DataDir:            "data",
		JWTExpiry:          24 * time.Hour,
		APIRateLimitRPM:    10,
		AccountWaitTimeout: 30 * time.Second,
	}
}
