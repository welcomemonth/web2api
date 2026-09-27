package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/web2api/internal/middleware"
)

func (app *App) registerAdminRouter(router *gin.Engine) {
	adminRouter := router.Group("/api/admin")
	adminRouter.Use(middleware.AuthAdmin)
	// TODO 这里使用中间拦截admin
	adminRouter.GET("/status", app.adminStatus)
}

func (app *App) adminStatus(c *gin.Context) {
	// TODO
	perAccount := []map[string]any{}

	for _, acc := range app.accounts.Snapshot() {
		perAccount = append(perAccount, map[string]any{
			"email": acc.Email, "status": acc.StatusCode, "inflight": acc.Inflight,
			"max_inflight": 1, "consecutive_failures": acc.ConsecutiveFailures,
			"rate_limit_strikes": acc.RateLimitStrikes, "last_request_finished": acc.LastRequestFinished,
		})
	}

	c.JSONP(http.StatusOK, map[string]any{
		"accounts":           app.accounts.Status(),
		"per_account":        "perAccount",
		"chat_id_pool":       "app.chatPool.Status()",
		"runtime":            map[string]any{"mode": "go", "goroutines_note": "not exposed"},
		"request_runtime":    map[string]any{"mode": "direct_http", "browser_required_for_requests": false, "description": "普通请求直连 HTTP，不经过浏览器"},
		"browser_automation": map[string]any{"mode": "playwright", "description": "Go 后端通过 Playwright 浏览器自动化支持邮箱激活"},
	})
}
