package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/restaurant/common/config"
	"github.com/SergeyBogomolovv/restaurant/common/db"
	"github.com/SergeyBogomolovv/restaurant/common/redis"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	db := db.MustConnect(cfg.PostgresURL)
	redis := redis.MustConnect(cfg.RedisURL)
	defer redis.Close()
	defer db.Close()
	logger := setupLogger(cfg.Env)

	slog.Info("starting application", slog.String("env", cfg.Env))

	app := app.New(logger, db, redis, cfg.Jwt, cfg.SSO.SecretKey)
	go app.Run(cfg.SSO.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()
	app.Shutdown()
	slog.Info("application stopped")
}

func setupLogger(env string) (logger *slog.Logger) {
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return
}
