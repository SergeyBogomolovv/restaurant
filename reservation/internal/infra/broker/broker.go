package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn *amqp.Connection
}

func NewRabbitMQBroker(conn *amqp.Connection) *RabbitMQBroker {
	return &RabbitMQBroker{conn: conn}
}

func (b *RabbitMQBroker) Publish(exchange, routingKey string, body []byte) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
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

	if err := ch.ExchangeDeclare("reservation_exchange", "topic", true, false, false, false, nil); err != nil {
		return err
	}

	return nil
}
