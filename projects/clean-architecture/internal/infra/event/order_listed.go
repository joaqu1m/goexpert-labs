package event

import "time"

type OrderListed struct {
	Name       string
	ReceivedAt time.Time
	Payload    any
}

func NewOrderListed() *OrderListed {
	return &OrderListed{
		Name: "OrderListed",
	}
}

func (e *OrderListed) GetName() string {
	return e.Name
}

func (e *OrderListed) GetReceivedAt() time.Time {
	return e.ReceivedAt
}

func (e *OrderListed) GetPayload() any {
	return e.Payload
}

func (e *OrderListed) SetPayload(payload any) {
	e.Payload = payload
}
