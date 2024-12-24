// mockgen -source=./internal/storage/store.go --destination=internal/mock/mock_store/mock_store.go -package=storage
package storage

import "github.com/wavinamayola/user-management/internal/models"

type Store interface {
	CreateUser(user models.UserRequest) (int, error)
	GetUser(id int) (models.User, error)
	UpdateUser(id int, user models.UserRequest) error
	DeleteUser(id int) error
}
