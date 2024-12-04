package broker

import (
	"github.com/SergeyBogomolovv/restaurant/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn *amqp.Connection
}

func NewRabbitMQBroker(conn *amqp.Connection) *RabbitMQBroker {
	return &RabbitMQBroker{conn: conn}
}

func (b *RabbitMQBroker) Publish(routingKey string, body []byte) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(constants.ReservationExchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (b *RabbitMQBroker) Setup() error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(constants.ReservationExchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	return nil
}
