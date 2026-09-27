package model

import "time"

// APIKeyStatus 表示 API Key 的启用状态。
type APIKeyStatus string

const (
	APIKeyStatusEnabled  APIKeyStatus = "enabled"
	APIKeyStatusDisabled APIKeyStatus = "disabled"
)

// APIKey 是 OpenAI-compatible 接口的鉴权凭据，对应 api_keys.json 中的一条记录。
// 第一版 Key 明文存储（仅本地自用，勿用于生产）。
type APIKey struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Key       string       `json:"key"`
	Status    APIKeyStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}
