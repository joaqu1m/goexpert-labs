package pkg

import "sync"

type EventInterface interface {
	GetName() string
	GetReceivedTimestamp() int64
	GetPayload() map[string]any
}

type EventHandlerInterface interface {
	HandleEvent(event EventInterface, wg *sync.WaitGroup) error
}

type EventDispatcherInterface interface {
	DispatchEvent(event EventInterface) error
	RegisterHandler(eventName string, handler EventHandlerInterface) error
	UnregisterHandler(eventName string, handler EventHandlerInterface) error
	GetHandlers(eventName string) ([]EventHandlerInterface, error)
}
