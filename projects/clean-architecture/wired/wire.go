//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

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

var setOrderRepositoryDependency = wire.NewSet(
	mysql.NewOrderRepository,
	wire.Bind(new(database.OrderRepositoryInterface), new(*mysql.OrderRepository)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setOrderListedEvent = wire.NewSet(
	event.NewOrderListed,
	wire.Bind(new(events.EventInterface), new(*event.OrderListed)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderListedEvent,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebCreateOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewOrderHandler,
	)
	return &web.OrderHandler{}
}

func NewWebListOrdersHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderListedEvent,
		web.NewOrderHandler,
	)
	return &web.OrderHandler{}
}
