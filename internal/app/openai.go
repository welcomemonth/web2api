package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/web2api/internal/services"
)

func (app *App) registerOpenAIRouter(engine *gin.Engine) {
	engine.GET("/v1/models", app.handleListModels)
}

func (app *App) handleListModels(c *gin.Context) {
	upstream, err := app.client.ListModelsFromPool(c.Request.Context())

	if err == nil && len(upstream) > 0 {
		c.JSON(http.StatusOK, services.BuildOpenAIModelList(upstream))
		return
	}

	c.JSON(http.StatusOK, services.BuildFallbackModelList())
}
