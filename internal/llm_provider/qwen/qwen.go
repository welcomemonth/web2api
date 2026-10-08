package qwen

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/model"
	"github.com/welcomemonth/web2api/internal/runtime"
	"github.com/welcomemonth/web2api/internal/utils"
)

// ==================== 配置 ====================
const (
	DefaultTimeoutMs = 30000.0 // 30s
	ShortTimeoutMs   = 5000.0  // 5s，用于快速失败
	RetryInterval    = 300 * time.Millisecond
	QwenBaseURL      = "https://chat.qwen.ai"
)

type Client struct {
	pool           *runtime.AccountPool
	mu             sync.Mutex
	browserContext playwright.BrowserContext
	cfg            *config.Config
}

type VerifyResult struct {
	Valid      bool
	StatusCode string
	Error      string
}

func NewClient(pool *runtime.AccountPool, browser playwright.Browser, cfg *config.Config) (*Client, error) {
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Locale: playwright.String("zh-CN"), //TODO 后续如果批量起账号，需要有这个 语言和位置等等
	})
	if err != nil {
		slog.Error("创建 context 失败", "error", err)
		return nil, err
	}
	client := &Client{
		browserContext: context,
		pool:           pool,
		cfg:            cfg,
	}
	return client, nil
}

func (q *Client) VerifyTokenDetail(ctx context.Context, token string) *VerifyResult {
	slog.Debug("verify Token Detail", "token", token)
	return nil
}

func (q *Client) VerifyAccountWithPwd(ctx context.Context, account *model.Account) *VerifyResult {
	var result = &VerifyResult{
		Valid: false,
	}
	// ============= 打开新页面 ==============
	page, err := q.browserContext.NewPage()
	if err != nil {
		result.Error = err.Error()
		result.StatusCode = ""
		return result
	}
	defer page.Close()
	// =============  打开目标网站 ===============
	if _, err = page.Goto(QwenBaseURL+"/auth", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		result.Error = err.Error()
		result.StatusCode = ""
		return result
	}
	// ============= 登陆 ===========
	if err = login(page, account.Email, account.Password); err != nil {
		result.Error = err.Error()
		result.StatusCode = ""
		return result
	}

	// ============== 等待登陆成功标识 ============
	userProfilelocator := page.Locator(`img[alt="User profile"]`)

	if err = userProfilelocator.WaitFor(playwright.LocatorWaitForOptions{
		// Timeout: &timeout, // 不需要传递，默认即30s超时即可
	}); err != nil {
		result.Error = err.Error()
		result.StatusCode = ""
		return result
	}

	// ========== 保存完整登录状态 ==========
	if _, err := q.browserContext.StorageState(playwright.BrowserContextStorageStateOptions{
		Path: playwright.String(q.cfg.DataDir + utils.GetEmailHashFilename(account.Email)),
	}); err != nil {
		result.Error = err.Error()
		result.StatusCode = ""
		return result
	}
	result.Valid = true
	return result
}

// ==================== 1. 清理遮挡物 ====================
// 尽量清除可能遮挡的悬浮层、Tooltip、Modal 残留
func ClearOverlays(page playwright.Page) {
	// 1. 按下ESC (关闭绝大多数非模态层)
	_ = page.Keyboard().Press("Escape")
	time.Sleep(300 * time.Millisecond)

	// 2. 可选：点击空白处关闭一些下拉/提示 (暂时注释，不确定点击哪个位置)
	// _ = page.Mouse().Click(10, 10)

	// 3. 可选：用 JS 强制移除常见遮罩（根据实际页面调整选择器）
	_, _ = page.Evaluate(`() => {
		document.querySelectorAll('.modal-backdrop, .overlay, .tooltip, [role="dialog"]').forEach(el => {
			if (el.style) el.style.display = 'none';
		});
	}`)
}

// ==================== 2. 安全点击（核心复用函数） ====================
// 多重策略：常规 → Force → JS 原生 click + 事件派发
// 最终失败返回 error
func SafeClick(locator playwright.Locator, timeoutMs float64) error {
	if timeoutMs <= 0 {
		timeoutMs = DefaultTimeoutMs
	}
	short := ShortTimeoutMs
	// 先等待元素出现
	if err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: &timeoutMs,
	}); err != nil {
		return fmt.Errorf("等待元素可见失败: %w", err)
	}
	// 策略 1: 常规点击（带 actionability 检查）
	err := locator.Click(playwright.LocatorClickOptions{
		Timeout: &short,
	})
	if err == nil {
		return nil
	}
	slog.Debug("常规点击失败，尝试 Force", "err", err)

	// 策略 2: Force 强制点击（忽略遮挡）
	err = locator.Click(playwright.LocatorClickOptions{
		Force:   playwright.Bool(true),
		Timeout: &short,
	})
	if err == nil {
		return nil
	}
	slog.Debug("Force 点击失败，尝试 JS 兜底", "err", err)
	// 策略 3: JS 终极兜底
	_, err = locator.Evaluate(`(el) => {
		el.scrollIntoView({block: 'center', behavior: 'instant'});
		el.click();
		el.dispatchEvent(new MouseEvent('click', {
			bubbles: true,
			cancelable: true,
			view: window
		}));
	}`, nil)
	if err != nil {
		return fmt.Errorf("所有点击策略均失败: %w", err)
	}
	return nil
}

// ==================== 3. 安全填充输入框 ====================
// Focus → Fill → JS 注入（触发 input/change）
func SafeFill(locator playwright.Locator, value string, timeoutMs float64) error {
	if timeoutMs <= 0 {
		timeoutMs = DefaultTimeoutMs
	}

	// 先尝试 Focus（比 Click 更轻量）
	err := locator.Focus(playwright.LocatorFocusOptions{Timeout: &timeoutMs})
	if err != nil {
		// Focus 失败就 Force Click
		if err2 := SafeClick(locator, timeoutMs); err2 != nil {
			return fmt.Errorf("聚焦/点击输入框失败: %w", err2)
		}
	}

	// 常规 Fill
	err = locator.Fill(value, playwright.LocatorFillOptions{Timeout: &timeoutMs})
	if err == nil {
		return nil
	}
	slog.Debug("常规 Fill 失败，使用 JS 注入", "err", err)

	// JS 强制注入 + 触发事件
	_, err = locator.Evaluate(`(el, val) => {
		el.focus();
		el.value = val;
		el.dispatchEvent(new Event('input',  { bubbles: true }));
		el.dispatchEvent(new Event('change', { bubbles: true }));
		// 部分框架还需要这个
		el.dispatchEvent(new Event('blur',   { bubbles: true }));
	}`, value)
	if err != nil {
		return fmt.Errorf("JS 注入值失败: %w", err)
	}
	return nil
}

// ==================== 4. 判断函数（可复用） ====================
// 元素是否可见且可交互
func IsInteractable(locator playwright.Locator, timeoutMs float64) bool {
	if timeoutMs <= 0 {
		timeoutMs = 3000.0
	}
	err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: &timeoutMs,
	})
	if err != nil {
		return false
	}
	// 额外检查 enabled
	enabled, err := locator.IsEnabled()
	return err == nil && enabled
}

// 元素是否存在（不要求可见）
func Exists(locator playwright.Locator, timeoutMs float64) bool {
	if timeoutMs <= 0 {
		timeoutMs = 2000.0
	}
	err := locator.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: &timeoutMs,
	})
	return err == nil
}

// 在填写账号密码之前，先判断并切换到密码登录模式
func switchToPasswordLoginIfNeeded(page playwright.Page) error {
	// 定位“使用密码登录”按钮（多重保险）
	passwordLoginBtn := page.Locator("div.qwenchat-email-otp-panel-return button").
		Or(page.GetByRole("button", playwright.PageGetByRoleOptions{
			Name:  "使用密码登录",
			Exact: playwright.Bool(true),
		})).
		Or(page.Locator("button:has-text('使用密码登录')"))

	// 判断按钮是否存在（给 3 秒超时即可）
	if !Exists(passwordLoginBtn, 3000) {
		slog.Info("未发现「使用密码登录」按钮，当前已是密码登录模式或页面结构不同")
		return nil
	}

	slog.Info("发现「使用密码登录」按钮，准备点击切换...")
	if err := SafeClick(passwordLoginBtn, DefaultTimeoutMs); err != nil {
		return fmt.Errorf("点击「使用密码登录」失败: %w", err)
	}

	// 点击后给页面一点时间切换表单
	time.Sleep(800 * time.Millisecond)
	slog.Info("✓ 已切换到密码登录模式")
	return nil
}

func login(page playwright.Page, account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("账号或密码不能为空")
	}

	// 0. 清理可能的遮挡
	ClearOverlays(page)

	// 1. 关键存在「使用密码登录」按钮，就先点击切换
	if err := switchToPasswordLoginIfNeeded(page); err != nil {
		return err
	}
	// 2. 填写邮箱
	emailInput := page.Locator("input[name='email']").
		Or(page.Locator("input[placeholder*='邮箱']")).
		Or(page.Locator(`[data-gtm-form-interact-field-id="0"]`))

	if err := SafeFill(emailInput, account, DefaultTimeoutMs); err != nil {
		return fmt.Errorf("填写邮箱失败: %w", err)
	}
	slog.Info("✓ 邮箱已填写")
	time.Sleep(800 * time.Millisecond)

	// 3. 填写密码
	passwordInput := page.Locator("input[name='password']").
		Or(page.Locator("input[type='password']")).
		Or(page.Locator(`[data-gtm-form-interact-field-id="1"]`))

	if err := SafeFill(passwordInput, password, DefaultTimeoutMs); err != nil {
		return fmt.Errorf("填写密码失败: %w", err)
	}
	slog.Info("✓ 密码已填写")
	time.Sleep(800 * time.Millisecond)

	// 登录按钮（优先按钮，失败再用 Enter）
	loginBtn := page.Locator("button[type='submit']").
		Or(page.GetByRole("button", playwright.PageGetByRoleOptions{
			Name:  "登录",
			Exact: playwright.Bool(true),
		})).
		Or(page.Locator("button:has-text('登录')"))

	if IsInteractable(loginBtn, ShortTimeoutMs) {
		if err := SafeClick(loginBtn, DefaultTimeoutMs); err != nil {
			slog.Warn("点击登录按钮失败，尝试回车", "err", err)
			if err := page.Keyboard().Press("Enter"); err != nil {
				return fmt.Errorf("回车提交也失败: %w", err)
			}
		}
	} else {
		// 按钮不可见时直接回车
		if err := page.Keyboard().Press("Enter"); err != nil {
			return fmt.Errorf("按下回车键失败: %v", err)
		}
	}

	slog.Info("✓ 已触发登录，等待页面响应...")
	return nil
}
