package usecase

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Consume(queue string, handler func(amqp.Delivery)) error
}

type paymentsUsecase struct {
	log    *slog.Logger
	broker Broker
}

func NewPaymentsUsecase(log *slog.Logger, broker Broker) *paymentsUsecase {
	return &paymentsUsecase{log: log, broker: broker}
}

func (u *paymentsUsecase) Run() {
	const op = "paymentsUsecase.Run"
	log := u.log.With(slog.String("op", op))
	log.Info("starting payments usecase")

	if err := u.broker.Consume("payments_reservation_queue", func(message amqp.Delivery) {
		log.Info("message received", "message", message.RoutingKey)
	}); err != nil {
		log.Error("failed to consume", "error", err)
	}
}
