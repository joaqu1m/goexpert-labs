package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetReceivedAt() time.Time
	GetPayload() any
	SetPayload(payload any)
}

type EventHandlerInterface interface {
	HandleEvent(event EventInterface, wg *sync.WaitGroup) error
}

type EventDispatcherInterface interface {
	DispatchEvent(event EventInterface) error
	RegisterHandler(eventName string, handler EventHandlerInterface) error
	UnregisterHandler(eventName string, handler EventHandlerInterface) error
	GetHandlers(eventName string) ([]EventHandlerInterface, error)
	ClearHandlers(eventName string) error
}
