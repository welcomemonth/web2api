package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/web2api/internal/app"
	"github.com/welcomemonth/web2api/internal/config"
	"github.com/welcomemonth/web2api/internal/pkg/logger"
)

func main() {
	cfg := config.Default()
	logger.Init(cfg.LogLevel, cfg.IsProd)
	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}
	// 1. 加载数据目录下的 JSON 到内存
	application, err := app.New(&cfg)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(-1)
	}

	backgroundCtx, stopBackground := context.WithCancel(context.Background())
	defer stopBackground()

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           application.Routes(),
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("server listen failed", "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("Web2API Go backend starting", "port", cfg.Port, "version", cfg.Version)
		errCh <- srv.Serve(ln)
	}()

	application.StartBackground(backgroundCtx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stop:
		slog.Info("shutdown server", "signal", sig.String())
		stopBackground()
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
		stopBackground()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Web2API Go backend stopped")
}
