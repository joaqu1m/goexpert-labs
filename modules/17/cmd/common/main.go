package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joaqu1m/goexpert-labs/modules/17/internal/db"
)

func main() {
	ctx := context.Background()

	conn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	queries := db.New(conn)

	err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:          uuid.New().String(),
		Name:        "Programming",
		Description: sql.NullString{String: "Learn to code", Valid: true},
	})
	if err != nil {
		panic(err)
	}

	categories, err := queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		println("Category ID:", category.ID)
		println("Name:", category.Name)
		if category.Description.Valid {
			println("Description:", category.Description.String)
		} else {
			println("Description: NULL")
		}
	}
}
