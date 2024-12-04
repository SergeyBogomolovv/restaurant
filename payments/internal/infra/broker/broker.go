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

func (b *RabbitMQBroker) Consume(queue string, handler func(amqp.Delivery)) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queue,
		"payments-service",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		handler(msg)
	}

	return nil
}

func (b *RabbitMQBroker) Setup() error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("payments_reservation_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(q.Name, "reservation.*", "reservation_exchange", false, nil); err != nil {
		return err
	}
	return nil
}
