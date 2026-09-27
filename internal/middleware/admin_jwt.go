// Package middleware 提供 HTTP 中间件，包括管理员 JWT 鉴权与后续的 API Key 鉴权、限流。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/web2api/internal/auth"
)

// ContextUsername 是管理员用户名在 gin.Context 中的键。
const ContextUsername = "admin_username"

// AdminJWT 校验 Authorization: Bearer <JWT>。
// 通过时把用户名写入 Context 供后续 Handler 使用；任何失败都返回 401。
func AdminJWT(jwtSvc *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			// 与「Token 无效」返回一致的响应，不区分缺失与格式错误。
			abortUnauthorized(c, "缺少或格式错误的 Authorization 头")
			return
		}

		claims, err := jwtSvc.Parse(token)
		if err != nil {
			abortUnauthorized(c, "Token 无效或已过期")
			return
		}

		c.Set(ContextUsername, claims.Username)
		c.Next()
	}
}

// bearerToken 从 Authorization 头中提取 Bearer Token。
// 头缺失、前缀不符或 Token 为空时返回 false。
func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}
