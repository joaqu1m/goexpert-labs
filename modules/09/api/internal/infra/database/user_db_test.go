package database

import (
	"testing"

	"github.com/joaqu1m/goexpert-labs/modules/09/api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser(
		"John Doe",
		"john.doe@gmail.com",
		"securepassword",
	)

	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var foundUser entity.User
	err = db.First(&foundUser, "email = ?", user.Email).Error
	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
}
