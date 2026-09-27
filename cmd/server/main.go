package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/web2api/internal/app"
	"github.com/welcomemonth/web2api/internal/auth"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/handler"
	"github.com/welcomemonth/web2api/internal/pkg/logger"
	"github.com/welcomemonth/web2api/internal/service"
)

func main() {
	cfg := config.Default()
	logger.Init(cfg.LogLevel, cfg.IsProd)

	// 1. 加载数据目录下的 JSON 到内存
	application, err := app.New(&cfg)
	if err != nil {
		slog.Error("加载数据失败", "content", err)
		os.Exit(-1)
	}

	// 2. 确保管理员存在（首次启动创建默认账号）
	adminSvc := service.NewAdminInitService(application.Admin)
	admin, err := adminSvc.Ensure()
	if err != nil {
		slog.Error("初始化管理员失败", "content", err)
		os.Exit(-1)
	}
	log.Printf("管理员已就绪: %s", admin.Username)

	// 3. 组装依赖
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry)
	authSvc := service.NewAdminAuthService(application.Admin, jwtSvc)

	// 4. 注册路由
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	handler.NewAdminHandler(authSvc).RegisterAdminRoutes(r)

	// 5. 启动
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("Web2API 启动", "地址", "http://localhost"+addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
