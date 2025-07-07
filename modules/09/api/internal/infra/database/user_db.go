package database

import (
	"github.com/joaqu1m/goexpert-labs/modules/09/api/internal/entity"
	pkg_entity "github.com/joaqu1m/goexpert-labs/modules/09/api/pkg/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (u *User) Create(user *entity.User) error {
	if err := u.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) FindByID(id pkg_entity.ID) (*entity.User, error) {
	var user entity.User
	if err := u.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
