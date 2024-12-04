package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/SergeyBogomolovv/restaurant/payments/internal/infra/broker"
	"github.com/SergeyBogomolovv/restaurant/payments/internal/usecase"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type application struct {
	log    *slog.Logger
	server *grpc.Server
}

func New(log *slog.Logger, amqpConn *amqp.Connection) *application {
	server := grpc.NewServer()

	broker := broker.NewRabbitMQBroker(amqpConn)
	if err := broker.Setup(); err != nil {
		panic(err)
	}

	usecase := usecase.NewPaymentsUsecase(log, broker)
	go usecase.Run()

	return &application{
		log:    log,
		server: server,
	}
}

func (a *application) Run(port int) {
	const op = "payments.Run"
	log := a.log.With(slog.String("op", op))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	log.Info("Starting gRPC server")
	if err := a.server.Serve(listener); err != nil {
		panic(err)
	}
}

func (a *application) Shutdown() {
	const op = "payments.Run"
	log := a.log.With(slog.String("op", op))

	a.server.GracefulStop()
	log.Info("gRPC server stopped")
}
