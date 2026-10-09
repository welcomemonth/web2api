package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mxschmitt/playwright-go"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/llm_provider/qwen"
	"github.com/welcomemonth/web2api/internal/runtime"
	"github.com/welcomemonth/web2api/internal/storage"
)

type App struct {
	Config         *config.Config
	engine         *gin.Engine
	client         *qwen.Client
	apiKeys        map[string]bool
	managedAPIKeys map[string]bool
	envAPIKeys     map[string]bool

	accounts      *runtime.AccountPool
	browser       playwright.Browser
	usersStore    *storage.JSONStore
	accountsStore *storage.JSONStore
}

// New 依据配置定位 DataDir 下的各 JSON 文件并加载到内存。
// 任一文件损坏都会返回错误，避免带病启动。
func New(cfg *config.Config) (*App, error) {
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.LogsDir, 0o755); err != nil {
		return nil, err
	}
	pw, err := playwright.Run() // todo app关闭的时候需要stop
	if err != nil {
		return nil, err
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // 设为 false 方便界面查看与手动操作
	})
	if err != nil {
		return nil, err
	}

	accountsStore := storage.NewJSONStore(cfg.DataDir+"/accounts.json", []any{})
	accounts := runtime.NewAccountPool(accountsStore, *cfg)
	qwenClient, err := qwen.NewClient(accounts, browser, cfg)
	if err != nil {
		return nil, err
	}
	a := &App{
		Config:         cfg,
		client:         qwenClient,
		accounts:       accounts,
		usersStore:     storage.NewJSONStore(cfg.DataDir+"/user.json", []any{}),
		accountsStore:  accountsStore,
		apiKeys:        make(map[string]bool),
		managedAPIKeys: make(map[string]bool),
		envAPIKeys:     make(map[string]bool),
	}
	if err := a.load(); err != nil {
		return nil, err
	}
	if err := a.accounts.Load(); err != nil {
		return nil, err
	}

	a.apiKeys, a.managedAPIKeys, a.envAPIKeys = loadAPIKeys(a.Config.DataDir + "apikeys.json")

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

func loadAPIKeys(path string) (map[string]bool, map[string]bool, map[string]bool) {
	managed := loadManagedAPIKeys(path)
	// envKeys := loadEnvAPIKeys()
	envKeys := make(map[string]bool)
	all := map[string]bool{}
	for key := range managed {
		all[key] = true
	}
	for key := range envKeys {
		all[key] = true
	}
	return all, managed, envKeys
}

func loadManagedAPIKeys(path string) map[string]bool {
	keys := map[string]bool{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return keys
	}
	var payload struct {
		Keys any `json:"keys"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		slog.Warn("failed to parse api_keys.json", "error", err)
		return keys
	}
	switch v := payload.Keys.(type) {
	case []any:
		for _, item := range v {
			if key := strings.TrimSpace(fmt.Sprint(item)); key != "" {
				keys[key] = true
			}
		}
	case string:
		for _, key := range splitEnvList(v) {
			keys[key] = true
		}
	}
	return keys
}

func splitEnvList(value string) []string {
	parts := regexp.MustCompile(`[,\s;]+`).Split(value, -1)
	out := []string{}
	seen := map[string]bool{}
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}
