package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/joaqu1m/goexpert-labs/modules/13/graph/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) CreateCategory(name string, description *string) (*model.Category, error) {
	id := uuid.New().String()
	query := "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3) RETURNING id, name, description"
	var category model.Category
	err := r.db.QueryRow(query, id, name, description).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
