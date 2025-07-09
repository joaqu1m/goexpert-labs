package pkg

import "sync"

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) DispatchEvent(event EventInterface) error {
	if handlers, exists := ed.handlers[event.GetName()]; exists {
		wg := &sync.WaitGroup{}
		wg.Add(len(handlers))
		for _, handler := range handlers {
			go handler.HandleEvent(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) RegisterHandler(eventName string, handler EventHandlerInterface) error {
	if _, exists := ed.handlers[eventName]; !exists {
		ed.handlers[eventName] = []EventHandlerInterface{}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) UnregisterHandler(eventName string, handler EventHandlerInterface) error {
	if handlers, exists := ed.handlers[eventName]; exists {
		for i, h := range handlers {
			if h == handler {
				ed.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *EventDispatcher) GetHandlers(eventName string) ([]EventHandlerInterface, error) {
	if handlers, exists := ed.handlers[eventName]; exists {
		return handlers, nil
	}
	return nil, nil
}
