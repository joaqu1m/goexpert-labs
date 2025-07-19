//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/joaqu1m/goexpert-labs/modules/19/product"
)

var setRepositoryDependency = wire.NewSet(
	product.NewProductRepository,
	wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository)),
)

func NewUsecase(db *sql.DB) *product.ProductUsecase {
	wire.Build(
		setRepositoryDependency,
		product.NewProductUsecase,
	)
	return &product.ProductUsecase{}
}
