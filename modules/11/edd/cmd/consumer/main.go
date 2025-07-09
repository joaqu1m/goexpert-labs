package main

import (
	"fmt"

	"github.com/joaqu1m/goexpert-labs/modules/11/edd/pkg"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, err := pkg.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs := make(chan amqp091.Delivery)

	go pkg.Consume(ch, msgs)

	fmt.Println("Waiting for messages...")

	for msg := range msgs {
		fmt.Println("Received message:", string(msg.Body))
		msg.Ack(false)
	}
}
