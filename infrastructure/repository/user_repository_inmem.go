package repository

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
	errors "github.com/pkg/errors"
)

type UserRepositoryInmem struct {
	table map[uuid.UUID]domain.User
}

func NewUserRepositoryInmem() *UserRepositoryInmem {
	return &UserRepositoryInmem{
		table: make(map[uuid.UUID]domain.User),
	}
}

func (db *UserRepositoryInmem) SaveUser(user domain.User) error {
	db.table[user.ID] = user
	return nil
}

func (db *UserRepositoryInmem) GetUserByID(id uuid.UUID) (*domain.User, error) {
	if user, ok := db.table[id]; ok {
		return &user, nil
	}
	return nil, errors.New("user not found")
}
