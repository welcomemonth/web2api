package app

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/runtime"
	"github.com/welcomemonth/web2api/internal/storage"
)

type App struct {
	Config *config.Config
	engine *gin.Engine

	accounts *runtime.AccountPool

	usersStore    *storage.JSONStore
	accountsStore *storage.JSONStore
}

// New 依据配置定位 DataDir 下的各 JSON 文件并加载到内存。
// 任一文件损坏都会返回错误，避免带病启动。
func New(cfg *config.Config) (*App, error) {
	a := &App{
		Config:        cfg,
		usersStore:    storage.NewJSONStore(cfg.DataDir+"/user.json", []any{}),
		accountsStore: storage.NewJSONStore(cfg.DataDir+"/accounts.json", []any{}),
	}
	if err := a.load(); err != nil {
		return nil, err
	}
	a.accounts = runtime.NewAccountPool(a.accountsStore, *a.Config)

	return a, nil
}

func (a *App) load() error {
	return nil
}

func (app *App) Routes() *gin.Engine {
	if app.engine != nil {
		return app.engine
	}
	app.engine = gin.Default()

	// engine.Use(middleware.RequestID(), middleware.GinLogger(), middleware.CORS())
	app.engine.MaxMultipartMemory = 50 << 20
	app.engine.Static("/static", "./uploads")

	app.registerAppHeartbeatRouter(app.engine)
	app.registerOpenAIRouter(app.engine)
	app.registerAnthropicRouter(app.engine)
	app.registerAdminRouter(app.engine)
	return app.engine
}

func (app *App) StartBackground(ctx context.Context) {
	if app == nil {
		return
	}
	// app.chatPool.Start(ctx)
	// if app.keepalive != nil {
	// 	app.keepalive.Start(ctx, app.keepaliveConfig())
	// }
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// app.cleanupContextArtifacts(ctx)
			}
		}
	}()
}
