package repository

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
	errors "github.com/pkg/errors"
)

type UserRepositoryFakeDB struct {
	table map[uuid.UUID]domain.User
}

func NewUserRepositoryFakeDB() *UserRepositoryFakeDB {
	return &UserRepositoryFakeDB{
		table: make(map[uuid.UUID]domain.User),
	}
}

func (db *UserRepositoryFakeDB) SaveUser(user domain.User) error {
	db.table[user.ID] = user
	return nil
}

func (db *UserRepositoryFakeDB) GetUserByID(id uuid.UUID) (*domain.User, error) {
	if user, ok := db.table[id]; ok {
		return &user, nil
	}
	return nil, errors.New("user not found")
}
