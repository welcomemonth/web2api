// Package handler 提供 HTTP 接口层，负责参数绑定、调用 Service 与错误转换。
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/web2api/internal/service"
)

// AdminHandler 处理管理员相关接口。
type AdminHandler struct {
	auth *service.AdminAuthService
}

// NewAdminHandler 创建管理员 Handler。
func NewAdminHandler(authSvc *service.AdminAuthService) *AdminHandler {
	return &AdminHandler{auth: authSvc}
}

// RegisterAdminRoutes 注册管理员公开路由（无需鉴权）。
func (h *AdminHandler) RegisterAdminRoutes(r gin.IRouter) {
	r.POST("/admin/login", h.Login)
}

// loginRequest 是 POST /admin/login 的请求体。
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login 处理管理员登录：校验参数、认证并返回 JWT。
func (h *AdminHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	token, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		// 签发失败等内部错误，不向外暴露细节。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
