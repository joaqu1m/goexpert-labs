package graph

import "github.com/joaqu1m/goexpert-labs/modules/13/internal/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CategoryRepository *database.CategoryRepository
	CourseRepository   *database.CourseRepository
}
