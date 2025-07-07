package database

import (
	"github.com/joaqu1m/goexpert-labs/modules/09/api/internal/entity"
	pkg_entity "github.com/joaqu1m/goexpert-labs/modules/09/api/pkg/entity"
)

type UserInterface interface {
	Create(user *entity.User) error
	FindByID(id pkg_entity.ID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
