package scheduler

import (
	"github.com/mxschmitt/playwright-go"

	"github.com/welcomemonth/web2api/internal/model"
)

// AccountRuntime 保存账号的运行时状态，不落盘。
// 由账号服务在创建/验证账号时构造，并注册到调度器。
type AccountRuntime struct {
	Account          *model.Account
	BrowserContext   playwright.BrowserContext
	Busy             bool
	CurrentRequestID string
}
