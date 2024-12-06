package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/restaurant/common/amqp"
	"github.com/SergeyBogomolovv/restaurant/common/config"
	"github.com/SergeyBogomolovv/restaurant/common/constants"
	"github.com/SergeyBogomolovv/restaurant/common/db"
	"github.com/SergeyBogomolovv/restaurant/common/redis"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/app"
)

func main() {
	cfg := config.MustLoad()
	db := db.MustConnect(cfg.PostgresURL)
	defer db.Close()

	redis := redis.MustConnect(cfg.RedisURL)
	defer redis.Close()

	amqpConn := amqp.MustConnect(cfg.AmqpURL)
	defer amqpConn.Close()

	logger := setupLogger(cfg.Env).With(slog.String("env", cfg.Env))

	app := app.New(logger, db, redis, amqpConn, cfg.Jwt, cfg.SSO.SecretKey)
	go app.Run(cfg.SSO.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()
	app.Shutdown()
}

func setupLogger(env string) (logger *slog.Logger) {
	switch env {
	case constants.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case constants.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case constants.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return
}
