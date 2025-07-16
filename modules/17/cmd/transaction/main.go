package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joaqu1m/goexpert-labs/modules/17/internal/db"
)

type CategoryDB struct {
	Conn    *sql.DB
	Queries *db.Queries
}

func NewCategoryDB(conn *sql.DB) *CategoryDB {
	return &CategoryDB{
		Conn:    conn,
		Queries: db.New(conn),
	}
}

func main() {
	ctx := context.Background()

	conn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	categoryDB := NewCategoryDB(conn)

	err = categoryDB.CreateCategories(ctx, []db.CreateCategoryParams{
		{
			ID:          uuid.New().String(),
			Name:        "Programming",
			Description: sql.NullString{String: "Learn to code", Valid: true},
		},
		{
			ID:          uuid.New().String(),
			Name:        "Databases",
			Description: sql.NullString{String: "Learn about SQL and NoSQL", Valid: true},
		},
	})
	if err != nil {
		panic(err)
	}

	categories, err := categoryDB.Queries.ListCategories(ctx)
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
	println("Categories created successfully!")
}

func (c *CategoryDB) CreateCategory(ctx context.Context, params db.CreateCategoryParams) error {
	return c.callTx(ctx, func(q *db.Queries) error {
		return q.CreateCategory(ctx, params)
	})
}

func (c *CategoryDB) CreateCategories(ctx context.Context, categories []db.CreateCategoryParams) error {
	return c.callTx(ctx, func(q *db.Queries) error {
		for _, category := range categories {
			if err := q.CreateCategory(ctx, category); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *CategoryDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := c.Queries.WithTx(tx)

	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit()
}
