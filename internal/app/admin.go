package app

import "github.com/gin-gonic/gin"

func (app *App) registerAdminRouter(router *gin.Engine) {
	adminRouter := router.Group("/api/admin")

	// TODO 这里使用中间拦截admin
	adminRouter.GET("/status", nil)
}
