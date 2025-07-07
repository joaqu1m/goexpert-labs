package main

import (
	"math/rand"
	"time"
)

type Payload struct {
	Message  string
	Provider string
}

func main() {

	ch1 := make(chan Payload)
	ch2 := make(chan Payload)
	ch3 := make(chan Payload)

	go produce("Kafka", ch1)
	go produce("SQS", ch2)
	go produce("RabbitMQ", ch3)

	for {
		select {
		case msg := <-ch1:
			println("Received from Kafka:", msg.Message)
		case msg := <-ch2:
			println("Received from SQS:", msg.Message)
		case msg := <-ch3:
			println("Received from RabbitMQ:", msg.Message)
		case <-time.After(3 * time.Second):
			println("Warning: No messages received within 3 seconds")
		}
	}
}

func produce(providerName string, ch chan<- Payload) {
	for {
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
		ch <- Payload{
			Message:  "Message from " + providerName,
			Provider: providerName,
		}
	}
}
