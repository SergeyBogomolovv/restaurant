package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/restaurant/common/amqp"
	"github.com/SergeyBogomolovv/restaurant/common/config"
	"github.com/SergeyBogomolovv/restaurant/common/db"
	"github.com/SergeyBogomolovv/restaurant/payments/internal/app"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	db := db.MustConnect(cfg.PostgresURL)
	defer db.Close()

	amqpConn := amqp.MustConnect(cfg.AmqpURL)
	defer amqpConn.Close()

	log := setupLogger(cfg.Env).With(slog.String("env", cfg.Env))

	app := app.New(log, amqpConn)
	go app.Run(cfg.Payments.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()
	app.Shutdown()
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
