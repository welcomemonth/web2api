package browser

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type BrowserManager struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
}

const initScript = `
		// console.log("[QWEN-SSE-CHUNK] 开始");
		const originalFetch = window.fetch;
		window.fetch = async function(...args) {
			// console.log("[QWEN-SSE-CHUNK] 检查到内容");
			const requestUrl = typeof args[0] === 'string' ? args[0] : args[0]?.url;
			
			// 发送原始请求
			const response = await originalFetch.apply(this, args);
			
			// 匹配 Qwen 目标的 SSE completions 接口
			if (requestUrl && requestUrl.includes('/api/v2/chat/completions')) {
				// 【关键】克隆一份响应流，否则 Qwen 前端读取时会因为 Stream Locked 报错卡死
				const cloneRes = response.clone();
				
				(async () => {
					try {
						const reader = cloneRes.body.getReader();
						const decoder = new TextDecoder('utf-8');
						while (true) {
							const { done, value } = await reader.read();
							if (done) break;
							
							// 将流片段解码并输出到 Console
							const chunkText = decoder.decode(value, { stream: true });
							console.log("[QWEN-SSE-CHUNK]", chunkText);
						}
					} catch (e) {
						console.error("[QWEN-SSE-ERROR]", e);
					}
				})();
			}
			
			return response;
		};
	`

// 初始化
func NewBrowserManager(endpoint string) (*BrowserManager, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("启动 Playwright 失败: %w", err)
	}

	// browser, err := pw.Firefox.Connect(endpoint)
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		pw.Stop() // 连接失败时记得关掉
		return nil, fmt.Errorf("连接 Camoufox 失败: %w", err)
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Locale:           playwright.String("zh-CN"),
		StorageStatePath: playwright.String("auth.json"),
	})
	if err != nil {
		// 这里根据需要决定要不要关 browser / pw
		return nil, fmt.Errorf("创建 Context 失败: %w", err)
	}

	if err := context.AddInitScript(playwright.Script{Content: playwright.String(initScript)}); err != nil {
		return nil, fmt.Errorf("注入 InitScript 失败: %v", err)
	}

	return &BrowserManager{
		pw:      pw,
		browser: browser,
		context: context,
	}, nil
}

// 获取 Context 供外部使用
func (m *BrowserManager) Context() playwright.BrowserContext {
	return m.context
}

// 程序结束时调用，统一清理
func (m *BrowserManager) Close() error {
	// 先关 context（可选，看你是否还要复用 browser）
	if m.context != nil {
		_ = m.context.Close()
	}

	// 注意：如果是 Connect 远程浏览器，通常不要 Close browser，
	// 否则可能把整个 Camoufox 关掉。根据你的实际需求决定。
	// if m.browser != nil {
	// 	_ = m.browser.Close()
	// }

	if m.pw != nil {
		return m.pw.Stop()
	}
	return nil
}
