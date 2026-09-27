package model

import "time"

// AccountStatus 表示账号的持久化状态。
type AccountStatus string

const (
	AccountStatusNormal   AccountStatus = "normal"
	AccountStatusError    AccountStatus = "error"
	AccountStatusDisabled AccountStatus = "disabled"
)

// Account 是 Qwen 账号的持久化模型，对应 accounts.json 中的一条记录。
// busy/idle 属于运行时状态，不放在这里，也不写入 JSON。
type Account struct {
	ID          string        `json:"id"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	AccessToken string        `json:"access_token"`
	Status      AccountStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}
