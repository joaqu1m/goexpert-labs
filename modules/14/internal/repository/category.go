package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/joaqu1m/goexpert-labs/modules/14/internal/pb"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) CreateCategory(name string, description *string) (*pb.Category, error) {
	id := uuid.New().String()
	query := "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3) RETURNING id, name, description"
	var category pb.Category
	err := r.db.QueryRow(query, id, name, description).Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetAllCategories() ([]*pb.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*pb.Category
	for rows.Next() {
		var category pb.Category
		if err := rows.Scan(&category.Id, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
