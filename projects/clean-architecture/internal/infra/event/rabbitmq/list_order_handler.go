package rabbitmq

import (
	"encoding/json"
	"sync"

	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
	"github.com/rabbitmq/amqp091-go"
)

type OrderListedHandler struct {
	RabbitMQChannel *amqp091.Channel
}

func NewOrderListedHandler(rabbitMQChannel *amqp091.Channel) *OrderListedHandler {
	return &OrderListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (o *OrderListedHandler) HandleEvent(event events.EventInterface, wg *sync.WaitGroup) error {
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
		"order.listed",
		false,
		false,
		message,
	)
}
