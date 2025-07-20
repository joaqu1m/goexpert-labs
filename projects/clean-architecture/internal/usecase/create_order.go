package usecase

import (
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/entity"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository   database.OrderRepositoryInterface
	EventDispatcher   events.EventDispatcherInterface
	OrderCreatedEvent events.EventInterface
}

func NewCreateOrderUseCase(
	OrderRepository database.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
	OrderCreatedEvent events.EventInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:   OrderRepository,
		EventDispatcher:   EventDispatcher,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	order.CalculateFinalPrice()

	if err := c.OrderRepository.Save(*order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreatedEvent.SetPayload(dto)

	if err := c.EventDispatcher.DispatchEvent(c.OrderCreatedEvent); err != nil {
		return OrderOutputDTO{}, err
	}

	return OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
