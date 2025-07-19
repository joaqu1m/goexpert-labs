package web

import (
	"encoding/json"
	"net/http"

	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/usecase"
	"github.com/joaqu1m/goexpert-labs/modules/20/pkg/events"
)

type OrderHandler struct {
	OrderRepository   database.OrderRepositoryInterface
	EventDispatcher   events.EventDispatcherInterface
	OrderCreatedEvent events.EventInterface
}

func NewOrderHandler(
	OrderRepository database.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
	OrderCreatedEvent events.EventInterface,
) *OrderHandler {
	return &OrderHandler{
		OrderRepository:   OrderRepository,
		EventDispatcher:   EventDispatcher,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.EventDispatcher, h.OrderCreatedEvent)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
