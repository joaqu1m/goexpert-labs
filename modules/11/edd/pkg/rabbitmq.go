package pkg

import (
	"github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp091.Channel, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func Consume(ch *amqp091.Channel, out chan<- amqp091.Delivery) {
	msgs, err := ch.Consume(
		"pedidos",
		"go-consumer",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		out <- msg
	}
}

func Publish(ch *amqp091.Channel, body []byte) error {
	err := ch.Publish(
		"amq.direct", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
