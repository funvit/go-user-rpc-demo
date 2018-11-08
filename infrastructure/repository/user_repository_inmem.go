package repository

import (
	"errors"

	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"
)

type UserRepositoryInmem struct {
	table map[uuid.UUID]domain.User
}

func NewUserRepositoryInmem() *UserRepositoryInmem {
	return &UserRepositoryInmem{
		table: make(map[uuid.UUID]domain.User),
	}
}
func (db *UserRepositoryInmem) AddUser(user domain.User) error {
	db.table[user.ID] = user
	return nil
}

func (db *UserRepositoryInmem) UpdateUser(user domain.User) (bool, error) {
	if _, ok := db.table[user.ID]; ok {
		db.table[user.ID] = user
		return true, nil
	}
	return false, errors.New("user not found")
}

func (db *UserRepositoryInmem) GetUserByID(id uuid.UUID) (*domain.User, error) {
	if user, ok := db.table[id]; ok {
		return &user, nil
	}
	// errors.New("user not found")
	return nil, nil
}
