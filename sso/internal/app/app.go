package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/SergeyBogomolovv/restaurant/common/config"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/handler"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/repo"
	"github.com/SergeyBogomolovv/restaurant/sso/internal/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type App struct {
	server *grpc.Server
	log    *slog.Logger
}

func New(log *slog.Logger, db *sqlx.DB, rdb *redis.Client, jwtConfig config.JwtConfig, secretKey string) *App {
	customerRepo := repo.NewCustomerRepo(db)
	adminRepo := repo.NewAdminRepo(db)
	waiterRepo := repo.NewWaiterRepo(db)

	tokensRepo := repo.NewTokensRepo(rdb, jwtConfig)
	authUsecase := usecase.NewAuthUsecase(log, customerRepo, waiterRepo, adminRepo, tokensRepo)
	registerUsecase := usecase.NewRegisterUsecase(log, customerRepo, waiterRepo, adminRepo, secretKey)

	server := grpc.NewServer()
	handler.RegisterGRPCHandler(server, authUsecase, registerUsecase)

	return &App{server: server, log: log}
}

func (a *App) Run(port int) {
	const op = "auth.Run"
	log := a.log.With(slog.String("op", op))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	log.Info("gRPC server started", "addr", listener.Addr().String())

	if err := a.server.Serve(listener); err != nil {
		panic(err)
	}
}

func (a *App) Shutdown() {
	const op = "auth.Shutdown"
	a.log.With(slog.String("op", op)).Info("stopping gRPC server")
	a.server.GracefulStop()
}
