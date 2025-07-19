package rabbitmq

import (
	"encoding/json"
	"sync"

	"github.com/joaqu1m/goexpert-labs/modules/20/pkg/events"
	"github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp091.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp091.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (o *OrderCreatedHandler) HandleEvent(event events.EventInterface, wg *sync.WaitGroup) error {
	defer wg.Done()

	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	message := amqp091.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	return o.RabbitMQChannel.Publish(
		"amq.direct",
		"order.created",
		false,
		false,
		message,
	)
}
