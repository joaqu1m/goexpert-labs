package database

import "github.com/joaqu1m/goexpert-labs/modules/20/internal/entity"

type OrderRepositoryInterface interface {
	Save(order entity.Order) error
}
