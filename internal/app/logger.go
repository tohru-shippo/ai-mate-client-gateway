package app

import (
	"log/slog"
	"os"
	"strings"
)

// 根据配置创建带服务基础字段的 slog logger。
func NewLogger(cfg *Config) *slog.Logger {
	level := parseLevel(cfg.Log.Level)
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if cfg.Log.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler).With(
		slog.String("service", cfg.AppName),
		slog.String("env", cfg.Env),
	)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
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
