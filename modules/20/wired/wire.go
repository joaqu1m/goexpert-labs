//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package wired

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/database"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/database/mysql"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/event"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/web"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/usecase"
	"github.com/joaqu1m/goexpert-labs/modules/20/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	mysql.NewOrderRepository,
	wire.Bind(new(database.OrderRepositoryInterface), new(*mysql.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.OrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewOrderHandler,
	)
	return &web.OrderHandler{}
}
