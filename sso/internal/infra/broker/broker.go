package broker

import (
	"encoding/json"

	"github.com/SergeyBogomolovv/restaurant/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn *amqp.Connection
}

func NewRabbitMQBroker(conn *amqp.Connection) *RabbitMQBroker {
	return &RabbitMQBroker{conn: conn}
}

func (b *RabbitMQBroker) Publish(routingKey string, payload any) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return ch.Publish(constants.SSOExchange, routingKey, false, false, amqp.Publishing{
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

	return ch.ExchangeDeclare(constants.SSOExchange, "topic", true, false, false, false, nil)
}
