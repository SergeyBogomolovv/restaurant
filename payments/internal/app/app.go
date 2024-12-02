package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type application struct {
	log    *slog.Logger
	server *grpc.Server
}

func New(log *slog.Logger) *application {
	server := grpc.NewServer()
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
