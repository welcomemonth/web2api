package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) registerAppHeartbeatRouter(router *gin.Engine) {
	router.GET("/api", app.handleAPI)
	router.GET("/healthz", app.handleHealth)
	router.GET("/readyz", app.handleReady)
	router.GET("/keepalive", app.handleKeepAlive)
	router.HEAD("/keepalive", app.handleKeepAlive)
}

func (app *App) handleAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Web2API Enterprise Gateway is running", "docs": "/docs", "version": app.Config.Version})
}

func (app *App) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (app *App) handleReady(c *gin.Context) {
	// TODO 这个返回账号状态还未实现
	c.JSON(http.StatusOK, gin.H{"status": "ready", "accounts": "app.accounts.Status()"})
}

func (app *App) handleKeepAlive(c *gin.Context) {
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": app.Config.Version})
}
