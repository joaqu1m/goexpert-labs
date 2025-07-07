package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser(
		"John Doe",
		"john.doe@gmail.com",
		"securepassword",
	)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@gmail.com", user.Email)
	assert.NotEmpty(t, user.Password)
}

func TestPasswordHashing(t *testing.T) {
	user, err := NewUser(
		"John Doe",
		"john.doe@gmail.com",
		"securepassword",
	)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "securepassword", user.Password, "password should be hashed")
	assert.True(t, user.ComparePassword("securepassword"))
}
