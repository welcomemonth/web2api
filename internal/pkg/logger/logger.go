package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init 根据配置构建 slog.Logger，并注册为全局默认 logger。
// 注册后其他包直接调用 slog.Info/Error 也会走同一个 Handler，
// 无需到处传递 logger 实例。
func Init(levelText string, isProd bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: ParseLevel(levelText),
	}

	var handler slog.Handler
	if isProd {
		handler = slog.NewJSONHandler(os.Stdout, opts) // 生产 JSON，便于日志采集
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts) // 开发可读
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}

// ParseLevel 解析日志级别文本，无法识别时回退到 Info。
func ParseLevel(levelText string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelText)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Fatal 记录一条错误日志后退出进程，补上 slog 缺失的 Fatal 语义。
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
