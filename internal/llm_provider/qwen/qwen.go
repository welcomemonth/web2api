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

const qwenBaseURL = "https://chat.qwen.ai"

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
	if _, err = page.Goto(qwenBaseURL+"/auth", playwright.PageGotoOptions{
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

func login(page playwright.Page, account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("账号或密码不能为空")
	}

	timeoutMs := 30000.0 // 30秒
	// ==========================================
	// 0. 预处理：尝试清除可能的遮挡物
	// ==========================================
	// 按 ESC 键可以关闭绝大多数非模态的悬浮窗、Tooltip 或残留的遮罩层
	page.Keyboard().Press("Escape")
	time.Sleep(500 * time.Millisecond) // 给页面 0.5 秒时间响应关闭动作
	// ==========================================
	// 1. 定位邮箱输入框 (多重保险策略)
	// ==========================================
	// 优先级 1: name="email" (最稳定，与后端表单提交强绑定，前端极少修改)
	// 优先级 2: placeholder 包含 "邮箱" (模糊匹配，即使改成"请输入邮箱"也能命中)
	// 优先级 3: data-gtm-form-interact-field-id="0" (埋点属性，为了数据统计连续性，通常常年不变)
	// ==========================================
	// 1. 填写邮箱 (跳过 Click，直接 Focus + Fill)
	// ==========================================
	emailInput := page.Locator("input[name='email']").
		Or(page.Locator("input[placeholder*='邮箱']"))

	// 策略 A: 尝试直接获取焦点 (比 Click 更轻量，不易被遮挡判定卡死)
	err := emailInput.Focus(playwright.LocatorFocusOptions{Timeout: &timeoutMs})
	if err != nil {
		fmt.Println("⚠️ 常规聚焦失败，可能被遮挡，尝试强制点击...")
		// 策略 B: 如果 Focus 失败，启用 Force 模式强制点击
		err = emailInput.Click(playwright.LocatorClickOptions{
			Force:   playwright.Bool(true), // 忽略遮挡检查，强制触发点击
			Timeout: &timeoutMs,
		})
		if err != nil {
			return fmt.Errorf("强制点击邮箱输入框失败: %v", err)
		}
	}

	// 策略 C: 填充内容。如果 Fill 依然被前端框架拦截，使用 JS 兜底
	err = emailInput.Fill(account)
	if err != nil {
		fmt.Println("⚠️ 常规填充失败，使用 JS 强制注入...")
		_, err = emailInput.Evaluate(`
			(el, val) => {
				el.value = val;
				el.focus();
				el.dispatchEvent(new Event('input', { bubbles: true }));
				el.dispatchEvent(new Event('change', { bubbles: true }));
			}
		`, account)
		if err != nil {
			return fmt.Errorf("JS 注入邮箱失败: %v", err)
		}
	}
	fmt.Println("✓ 邮箱已填写")
	time.Sleep(2000 * time.Millisecond)
	// ==========================================
	// 2. 定位密码输入框 (多重保险策略)
	// ==========================================
	// 优先级 1: name="password" (最稳定)
	// 优先级 2: type="password" (密码框的通用特征)
	// 优先级 3: data-gtm-form-interact-field-id="1" (埋点属性)
	passwordInput := page.Locator("input[name='password']").
		Or(page.Locator("input[type='password']"))

	err = passwordInput.Focus(playwright.LocatorFocusOptions{Timeout: &timeoutMs})
	if err != nil {
		err = passwordInput.Click(playwright.LocatorClickOptions{
			Force:   playwright.Bool(true),
			Timeout: &timeoutMs,
		})
		if err != nil {
			return fmt.Errorf("强制点击密码输入框失败: %v", err)
		}
	}

	err = passwordInput.Fill(password)
	if err != nil {
		_, err = passwordInput.Evaluate(`
			(el, val) => {
				el.value = val;
				el.focus();
				el.dispatchEvent(new Event('input', { bubbles: true }));
				el.dispatchEvent(new Event('change', { bubbles: true }));
			}
		`, password)
		if err != nil {
			return fmt.Errorf("JS 注入密码失败: %v", err)
		}
	}
	fmt.Println("✓ 密码已填写")
	time.Sleep(2000 * time.Millisecond)

	// ==========================================
	// 3. 点击登录按钮
	// ==========================================
	slog.Debug("尝试通过键盘回车键提交表单...")
	err = page.Keyboard().Press("Enter")
	if err != nil {
		return fmt.Errorf("按下回车键失败: %v", err)
	}
	// loginBtn := page.Locator("button[type='submit']").
	// 	Or(page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "登录", Exact: playwright.Bool(true)}))

	// err = loginBtn.WaitFor(playwright.LocatorWaitForOptions{Timeout: &timeoutMs})
	// if err != nil {
	// 	return fmt.Errorf("等待登录按钮失败: %v", err)
	// }
	// // 策略 1: 尝试常规点击 (设置较短超时，避免干等 30 秒)
	// shortTimeout := 5000.0 // 5秒
	// err = loginBtn.Click(playwright.LocatorClickOptions{
	// 	Timeout: &shortTimeout,
	// })
	// if err != nil {
	// 	fmt.Println("⚠️ 常规点击被遮挡或超时，尝试策略 2：Force 强制点击...")

	// 	// 策略 2: 启用 Force 模式，告诉 Playwright 忽略所有遮挡检查，直接派发点击事件
	// 	err = loginBtn.Click(playwright.LocatorClickOptions{
	// 		Force:   playwright.Bool(true),
	// 		Timeout: &shortTimeout,
	// 	})

	// 	if err != nil {
	// 		fmt.Println("⚠️ Force 点击仍失败，尝试策略 3：JavaScript 终极兜底点击...")

	// 		// 策略 3: 完全绕过 Playwright 的模拟鼠标机制，直接在浏览器 DOM 层面执行点击
	// 		_, err = loginBtn.Evaluate(`
	// 			(el) => {
	// 				// 尝试原生点击
	// 				el.click();
	// 				// 如果原生点击被前端框架拦截，强制派发一个冒泡的 MouseEvent
	// 				el.dispatchEvent(new MouseEvent('click', {
	// 					bubbles: true,
	// 					cancelable: true,
	// 					view: window
	// 				}));
	// 			}
	// 		`, nil)
	// 		if err != nil {
	// 			return fmt.Errorf("所有点击策略均失败，可能存在严重的 DOM 遮挡或前端拦截: %v", err)
	// 		}
	// 	}
	// }
	slog.Info("✓ 已点击登录按钮，等待页面响应...")

	return nil
}
