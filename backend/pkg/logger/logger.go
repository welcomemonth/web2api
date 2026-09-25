package logger

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(levelText string) *slog.Logger {
	level := slog.LevelInfo
	switch normalizeLower(levelText) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func normalizeLower(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
