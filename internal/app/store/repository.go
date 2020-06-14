package store

import "github.com/assanoff/http-crud-server/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	GetUserByID(int) (*model.User, error)
	GetUserByField(string, string) (*model.User, error)
	GetUsers() ([]*model.User, error)
	UpdateUserByID(int, *model.User) (*model.User, error)
	DeleteUserByID(id int) error
}
