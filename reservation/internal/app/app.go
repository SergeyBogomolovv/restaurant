package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/SergeyBogomolovv/restaurant/reservation/internal/handler"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/infra/broker"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/repo"
	"github.com/SergeyBogomolovv/restaurant/reservation/internal/usecase"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type App struct {
	server     *grpc.Server
	log        *slog.Logger
	stopTicker context.CancelFunc
}

func New(log *slog.Logger, db *sqlx.DB, amqpConn *amqp.Connection) *App {
	server := grpc.NewServer()

	repo := repo.NewReservationRepo(db)

	broker := broker.NewRabbitMQBroker(amqpConn)
	if err := broker.Setup(); err != nil {
		panic(err)
	}

	usecase := usecase.NewReservationUsecase(log, repo, broker)
	ctx, cancel := context.WithCancel(context.Background())
	go usecase.RunEndedReservationsChecker(ctx, time.Hour)

	handler.RegisterGRPCHandler(server, usecase)

	return &App{server: server, log: log, stopTicker: cancel}
}

func (a *App) Run(port int) {
	const op = "reservation.Run"
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
	const op = "reservation.Shutdown"
	log := a.log.With(slog.String("op", op))
	a.server.GracefulStop()
	a.stopTicker()
	log.Info("gRPC server stopped")
}
