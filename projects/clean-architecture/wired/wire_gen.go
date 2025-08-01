// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wired

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/database/mysql"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/event"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/web"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/usecase"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
)

// Injectors from wire.go:

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	orderRepository := mysql.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, eventDispatcher, orderCreated)
	return createOrderUseCase
}

func NewListOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrdersUseCase {
	orderRepository := mysql.NewOrderRepository(db)
	orderListed := event.NewOrderListed()
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository, eventDispatcher, orderListed)
	return listOrdersUseCase
}

func NewWebCreateOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	orderRepository := mysql.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	orderHandler := web.NewOrderHandler(orderRepository, eventDispatcher, orderCreated)
	return orderHandler
}

func NewWebListOrdersHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	orderRepository := mysql.NewOrderRepository(db)
	orderListed := event.NewOrderListed()
	orderHandler := web.NewOrderHandler(orderRepository, eventDispatcher, orderListed)
	return orderHandler
}

// wire.go:

var setOrderRepositoryDependency = wire.NewSet(mysql.NewOrderRepository, wire.Bind(new(database.OrderRepositoryInterface), new(*mysql.OrderRepository)))

var setOrderCreatedEvent = wire.NewSet(event.NewOrderCreated, wire.Bind(new(events.EventInterface), new(*event.OrderCreated)))

var setOrderListedEvent = wire.NewSet(event.NewOrderListed, wire.Bind(new(events.EventInterface), new(*event.OrderListed)))
