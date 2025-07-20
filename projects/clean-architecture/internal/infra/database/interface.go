package database

import "github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/entity"

type OrderRepositoryInterface interface {
	List() ([]entity.Order, error)
	Save(order entity.Order) error
}
