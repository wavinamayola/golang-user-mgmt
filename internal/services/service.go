package service

import (
	"github.com/wavinamayola/user-management/internal/services/user"
	"github.com/wavinamayola/user-management/internal/storage"
)

type Service struct {
	UserService *user.User
}

func NewAPI(store *storage.Storage) *Service {
	u := user.New(store)

	return &Service{
		UserService: u,
	}
}
