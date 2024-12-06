package usecase

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentsBroker interface {
	Consume(queue string, handler func(amqp.Delivery)) error
}

type paymentsUsecase struct {
	log    *slog.Logger
	broker PaymentsBroker
}

func NewPaymentsUsecase(log *slog.Logger, broker PaymentsBroker) *paymentsUsecase {
	return &paymentsUsecase{log: log, broker: broker}
}

func (u *paymentsUsecase) Run() {
	const op = "paymentsUsecase.Run"
	log := u.log.With(slog.String("op", op))
	log.Info("starting payments usecase")

	if err := u.broker.Consume("payments.reservation_created_queue", func(message amqp.Delivery) {
		log.Info("message received", "message", message.RoutingKey)
		message.Ack(false)
	}); err != nil {
		log.Error("failed to consume", "error", err)
	}
}
