package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/tohru-shippo/ai-mate-client-gateway/docs"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/app"
)

// @title           AI Mate Client Gateway API
// @version         1.0
// @description     用户端 API 网关（BFF），将 HTTP 请求翻译为 ai-mate-server gRPC 调用。
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     Bearer 访问令牌，格式 "Bearer <token>"
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := app.Load()
	if err != nil {
		slog.Error("load config failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger := app.NewLogger(cfg)
	slog.SetDefault(logger)

	shutdownTracer, err := app.InitTracer(ctx, cfg, logger)
	if err != nil {
		logger.Error("init tracer failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			logger.Error("shutdown tracer failed", slog.Any("error", err))
		}
	}()

	httpServer, err := app.NewServer(cfg, logger)
	if err != nil {
		logger.Error("build http server failed", slog.Any("error", err))
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("client gateway starting", slog.String("addr", cfg.HTTP.Addr), slog.String("env", cfg.Env))
		errCh <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("client gateway stopped")
}
