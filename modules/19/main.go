package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	usecase := NewUsecase(db)

	product, err := usecase.GetProductByID(1)
	if err != nil {
		panic(err)
	}
	println("Product ID:", product.ID)
	println("Product Name:", product.Name)
	println("Product retrieved successfully")
}
