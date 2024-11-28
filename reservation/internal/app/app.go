package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	server *grpc.Server
	log    *slog.Logger
}

func New(log *slog.Logger, db *sqlx.DB) *App {
	//TODO: add deps
	server := grpc.NewServer()

	return &App{server: server, log: log}
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
	a.log.With(slog.String("op", op)).Info("stopping gRPC server")
	a.server.GracefulStop()
}
