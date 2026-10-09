package app

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/web2api/internal/middleware"
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/utils"
)

func (app *App) registerAdminRouter(router *gin.Engine) {
	adminRouter := router.Group("/api/admin")
	adminRouter.Use(middleware.AuthAdmin)
	// TODO 这里使用中间拦截admin
	adminRouter.GET("/status", app.adminStatus)
	adminRouter.GET("/accounts", app.adminListAccounts)
	adminRouter.POST("/accounts", app.adminAddAccount)
	adminRouter.POST("/accounts/:email/verify", app.adminVerifyAccount)
	adminRouter.DELETE("/accounts/:email", app.adminDeleteAccount)

	adminRouter.GET("/keys", app.adminGetKeys)
	adminRouter.POST("/keys", app.adminCreateKeys)
	adminRouter.DELETE("/keys/:key", app.adminDeleteKey)

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
		"accounts":    app.accounts.Status(),
		"per_account": perAccount,
		// "chat_id_pool":       app.chatPool.Status(), // 对话框预热池？是不是不需要了
		"runtime":            map[string]any{"mode": "go", "goroutines_note": "not exposed"},
		"request_runtime":    map[string]any{"mode": "direct_http", "browser_required_for_requests": false, "description": "普通请求直连 HTTP，不经过浏览器"},
		"browser_automation": map[string]any{"mode": "playwright", "description": "Go 后端通过 Playwright 浏览器自动化支持邮箱激活"},
	})
}

func (app *App) adminListAccounts(c *gin.Context) {

	accounts := []map[string]any{}

	for _, acc := range app.accounts.Snapshot() {
		accounts = append(accounts, map[string]any{
			"email": acc.Email, "password": acc.Password, "token": acc.Token, "cookies": acc.Cookies,
			"username": acc.Username, "activation_pending": acc.ActivationPending, "status_code": acc.StatusCode,
			"source": acc.Source, "env_name": acc.EnvName,
			"last_error": acc.LastError, "last_request_started": acc.LastRequestStarted, "last_request_finished": acc.LastRequestFinished,
			"consecutive_failures": acc.ConsecutiveFailures, "rate_limit_strikes": acc.RateLimitStrikes,
			"valid": acc.Valid, "inflight": acc.Inflight, "rate_limited_until": acc.RateLimitedUntil,
			"rate_limits": acc.RateLimits, // TODO 原来有一个cloneRateLimits
		})
	}
	c.JSONP(http.StatusOK, map[string]any{"accounts": accounts})
}

func (app *App) adminAddAccount(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "detail": "Invalid JSON body"})
		return
	}
	// 由于这个body不是结构体，因此解析之后，也只好自己手动判断需要字段是否存在
	token := utils.StringValue(body, "token", "")
	email := utils.StringValue(body, "email", "")
	if token == "" && email == "" {
		c.JSON(http.StatusBadRequest, map[string]any{"detail": "At least one of token or account is required"})
		return
	}

	acc := model.Account{
		Email:      utils.StringValue(body, "email", fmt.Sprintf("manual_%d@qwen", time.Now().Unix())),
		Password:   utils.StringValue(body, "password", ""),
		Token:      token,
		Cookies:    utils.StringValue(body, "cookies", ""),
		Username:   utils.StringValue(body, "username", ""),
		StatusCode: "valid",
	}
	slog.Info("new account", "account", acc)

	verify := app.client.VerifyAccountWithPwd(c.Request.Context(), &acc)
	if !verify.Valid {
		c.JSON(http.StatusOK, map[string]any{"ok": false, "error": "Invalid account (验证失败，请确认账号有效或密码正确)", "status_code": verify.StatusCode, "detail": verify.Error})
		return
	}
	if err := app.accounts.Add(acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "email": acc.Email})

}

func (app *App) adminDeleteAccount(c *gin.Context) {
	email := c.Param("email")

	if !utils.VerifyEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if err := app.accounts.Remove(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func (app *App) adminVerifyAccount(c *gin.Context) {
	email := c.Param("email")
	if !utils.VerifyEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	for _, acc := range app.accounts.Snapshot() {
		if acc.Email == email {
			verify := app.client.VerifyTokenDetail(c.Request.Context(), acc.Token)
			// _ = app.accounts.MarkVerification(acc.Email, verify)
			c.JSON(http.StatusOK, gin.H{"email": email, "valid": verify.Valid, "status_code": verify.StatusCode, "error": verify.Error, "refreshed": false})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"detail": "Account not found",
	})

}

func (app *App) adminGetKeys(c *gin.Context) {
	keys := []string{}
	items := []map[string]any{}

	for key := range app.apiKeys {
		keys = append(keys, key)
		source := "managed"
		label := "面板创建 Key"
		if app.envAPIKeys[key] {
			source = "env"
			label = "环境变量注入 Key"
		}
		items = append(items, map[string]any{"key": key, "source": source, "label": label})
	}
	sort.Strings(keys)
	sort.Slice(items, func(i, j int) bool {
		return fmt.Sprint(items[i]["key"]) < fmt.Sprint(items[j]["key"])
	})
	c.JSON(http.StatusOK, gin.H{
		"keys":  keys,
		"items": items,
	})
}

func (app *App) adminCreateKeys(c *gin.Context) {
	var body struct {
		Mode string `json:"mode"`
		Key  string `json:"key"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}
	mode := utils.NormalizeLower(body.Mode)
	if mode == "" {
		mode = "auto"
	}
	key := strings.TrimSpace(body.Key)

	if mode == "custom" {
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "自定义 Key 不能为空",
			})
			return
		}
		if strings.ContainsAny(key, " \t\r\n") {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "自定义 Key 不能包含空白字符",
			})
			return
		}
	} else {
		buf := make([]byte, 24)
		_, _ = cryptorand.Read(buf)
		key = "sk-" + hex.EncodeToString(buf)
	}

	if app.apiKeys[key] {
		c.JSON(http.StatusConflict, gin.H{
			"detail": "API Key 已存在",
		})
		return
	}

	app.apiKeys[key] = true
	app.managedAPIKeys[key] = true
	_ = saveAPIKeys(app.Config.DataDir+"apikeys.json", app.managedAPIKeys)
	c.JSON(http.StatusOK, gin.H{"ok": true, "key": key})
}

func saveAPIKeys(path string, keys map[string]bool) error {
	list := make([]string, 0, len(keys))
	for k := range keys {
		list = append(list, k)
	}
	sort.Strings(list)
	return utils.WriteJSONFileLocked(path, map[string]any{"keys": list})
}

func (app *App) adminDeleteKey(c *gin.Context) {
	key := c.Param("key")
	if app.envAPIKeys[key] {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "环境变量注入 Key 不能在面板删除，请移除对应环境变量后重启服务"})
		return
	}
	delete(app.apiKeys, key)
	delete(app.managedAPIKeys, key)
	_ = saveAPIKeys(app.Config.DataDir+"apikeys.json", app.managedAPIKeys)
	c.JSON(http.StatusOK, gin.H{"ok": true, "key": key})
}
