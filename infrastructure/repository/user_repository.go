package repository

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	AddUser(user domain.User) error
	UpdateUser(user domain.User) (bool, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
}
