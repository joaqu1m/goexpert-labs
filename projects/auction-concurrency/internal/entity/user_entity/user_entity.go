package user_entity

import (
	"context"

	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
