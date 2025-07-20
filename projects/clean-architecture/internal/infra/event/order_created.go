package event

import "time"

type OrderCreated struct {
	Name       string
	ReceivedAt time.Time
	Payload    any
}

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{
		Name: "OrderCreated",
	}
}

func (e *OrderCreated) GetName() string {
	return e.Name
}

func (e *OrderCreated) GetReceivedAt() time.Time {
	return e.ReceivedAt
}

func (e *OrderCreated) GetPayload() any {
	return e.Payload
}

func (e *OrderCreated) SetPayload(payload any) {
	e.Payload = payload
}
