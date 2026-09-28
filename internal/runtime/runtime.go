package runtime

import (
	"log/slog"
	"sort"
	"sync"

	"github.com/mxschmitt/playwright-go"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/storage"
	"github.com/welcomemonth/web2api/internal/utils"
)

const AccountReadySetThreshold = 1

// AccountPool 保存账号的运行时状态，不落盘。
// 由账号服务在创建/验证账号时构造，并注册到调度器。
type AccountPool struct {
	mu         sync.Mutex
	settings   config.Config
	store      *storage.JSONStore
	accounts   []*model.Account
	webClients map[string]playwright.Page

	maxInflightPerAccount  int
	globalInUse            int
	globalMaxInflight      int
	recommendedConcurrency int
	maxQueueSize           int
	readySetEnabled        bool
}

func NewAccountPool(store *storage.JSONStore, settings config.Config) *AccountPool {
	return &AccountPool{
		store:                 store,
		settings:              settings,
		maxInflightPerAccount: 1,
	}
}

func (p *AccountPool) Load() error {
	if err := p.store.Ensure(); err != nil {
		return err
	}
	var data []model.Account
	if err := p.store.LoadInto(&data); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.accounts = make([]*model.Account, 0, len(data))
	for i := range data {
		data[i].Normalize()
		data[i].MigrateLegacyRateLimit()
		p.accounts = append(p.accounts, &data[i])
	}
	// 不从环境变量中读取账号
	// for _, envAcc := range loadEnvAccounts() {
	// 	envAcc.normalize()
	// 	envAcc.migrateLegacyRateLimit(p.settings)
	// 	replaced := false
	// 	for i, existing := range p.accounts {
	// 		if existing.Email == envAcc.Email && envAcc.Email != "" {
	// 			cp := envAcc
	// 			p.accounts[i] = &cp
	// 			replaced = true
	// 			break
	// 		}
	// 	}
	// 	if !replaced {
	// 		cp := envAcc
	// 		p.accounts = append(p.accounts, &cp)
	// 	}
	// }
	p.resetLocked()
	slog.Info("loaded upstream accounts", "count", len(p.accounts))
	return nil
}

func (p *AccountPool) resetLocked() {
	valid := 0
	available := 0
	for _, acc := range p.accounts {
		if acc.Valid {
			valid++
		}
		if acc.AvailableFor("chat") {
			available++
		}
	}
	p.recommendedConcurrency = available * p.maxInflightPerAccount
	p.globalMaxInflight = p.recommendedConcurrency
	p.maxQueueSize = p.recommendedConcurrency
	p.readySetEnabled = valid >= utils.MaxInt(AccountReadySetThreshold, 1)
}

func (p *AccountPool) Status() map[string]any {
	p.mu.Lock()
	defer p.mu.Unlock()
	total := len(p.accounts)
	valid := 0
	available := 0
	availableImage := 0
	availableVideo := 0
	for _, acc := range p.accounts {
		if acc.Valid {
			valid++
		}
		if acc.AvailableFor("chat") && acc.Inflight < p.maxInflightPerAccount {
			available++
		}
		if acc.AvailableFor("image") && acc.Inflight < p.maxInflightPerAccount {
			availableImage++
		}
		if acc.AvailableFor("video") && acc.Inflight < p.maxInflightPerAccount {
			availableVideo++
		}
	}
	return map[string]any{
		"total":                    total,
		"valid":                    valid,
		"available":                available,
		"available_chat":           available,
		"available_image":          availableImage,
		"available_video":          availableVideo,
		"max_inflight_per_account": p.maxInflightPerAccount,
		"recommended_concurrency":  p.recommendedConcurrency,
		"global_max_inflight":      p.globalMaxInflight,
		"max_queue_size":           p.maxQueueSize,
		"global_in_use":            p.globalInUse,
		"ready_set_enabled":        p.readySetEnabled,
		"ready_set_threshold":      AccountReadySetThreshold,
	}
}

func (p *AccountPool) Snapshot() []model.Account {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]model.Account, 0, len(p.accounts))
	if len(p.accounts) < 1 {
		return out
	}
	for _, acc := range p.accounts {
		// acc.syncLegacyRateLimit() // TODO 需要深拷贝其他变量
		cp := *acc
		// cp.RateLimits = cloneRateLimits(acc.RateLimits)
		// cp.StatusCode = acc.status()
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Email < out[j].Email })
	return out
}
