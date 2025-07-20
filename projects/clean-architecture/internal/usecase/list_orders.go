package usecase

import (
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
)

type ListOrdersUseCase struct {
	OrderRepository  database.OrderRepositoryInterface
	EventDispatcher  events.EventDispatcherInterface
	OrderListedEvent events.EventInterface
}

func NewListOrdersUseCase(
	OrderRepository database.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
	OrderListedEvent events.EventInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository:  OrderRepository,
		EventDispatcher:  EventDispatcher,
		OrderListedEvent: OrderListedEvent,
	}
}

func (c *ListOrdersUseCase) Execute(input OrderInputDTO) ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	var output []OrderOutputDTO
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		output = append(output, dto)
	}

	c.OrderListedEvent.SetPayload(output)
	if err := c.EventDispatcher.DispatchEvent(c.OrderListedEvent); err != nil {
		return nil, err
	}

	return output, nil
}
