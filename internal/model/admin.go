package model

// Admin 管理员账号。第一版仅支持单个管理员，密码按约定明文存储（仅本地自用，勿用于生产）。
type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
