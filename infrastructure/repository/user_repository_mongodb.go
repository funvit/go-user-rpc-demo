package repository

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/google/uuid"

	errors "github.com/pkg/errors"
)

// UserRepositoryMongodb is struct of params
type UserRepositoryMongodb struct {
	//FIXME: добавить conn
}

// NewUserRepositoryMongodb create new instance of UserRepositoryMongodb
func NewUserRepositoryMongodb() *UserRepositoryMongodb {
	return &UserRepositoryMongodb{}
}

const repoName = "mongodb"

// SaveUser implements user adding to DB
func (db *UserRepositoryMongodb) SaveUser(user domain.User) error {
	return errors.Wrap(
		errors.New("not implemented"),
		repoName)
}

//GetUserByID implements get user from DB
func (db *UserRepositoryMongodb) GetUserByID(ID uuid.UUID) (*domain.User, error) {
	return nil, errors.Wrap(
		errors.New("not implemented"),
		repoName)
}
