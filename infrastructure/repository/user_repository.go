package repository

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	SaveUser(user domain.User) error
	GetUserByID(id uuid.UUID) (*domain.User, error)
}
