package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/wavinamayola/user-management/internal/models"
)

var validate = validator.New()

type Store interface {
	CreateUser(user models.UserRequest) (int, error)
	GetUser(id int) (models.User, error)
	UpdateUser(id int, user models.UserRequest) error
	DeleteUser(id int) error
}

type User struct {
	store Store
}

func New(s Store) *User {
	return &User{
		store: s,
	}
}
