package web

import (
	"encoding/json"
	"net/http"

	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/usecase"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
)

type OrderHandler struct {
	OrderRepository database.OrderRepositoryInterface
	EventDispatcher events.EventDispatcherInterface
	OrderEvent      events.EventInterface
}

func NewOrderHandler(
	OrderRepository database.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
	OrderEvent events.EventInterface,
) *OrderHandler {
	return &OrderHandler{
		OrderRepository: OrderRepository,
		EventDispatcher: EventDispatcher,
		OrderEvent:      OrderEvent,
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	listOrders := usecase.NewListOrdersUseCase(h.OrderRepository, h.EventDispatcher, h.OrderEvent)
	orders, err := listOrders.Execute(usecase.OrderInputDTO{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.EventDispatcher, h.OrderEvent)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
