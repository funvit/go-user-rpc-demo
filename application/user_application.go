package application

import (
	"github.com/funvit/go-user-rpc-demo/domain"
	"github.com/funvit/go-user-rpc-demo/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Context struct {
	UserRepo repository.UserRepository // FIXME: как то переделать на interface?
}

func (c *Context) AddUser(login string) (*domain.User, error) {
	logrus.Debugf("Application: adding user '%s'", login)
	user := domain.CreateNewUser(login)
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := c.UserRepo.AddUser(*user); err != nil {
		return nil, err
	} else {
		logrus.Debugf("Application: added user '%s', new ID=%s", user.Login, user.ID)
		return user, nil
	}
}

func (c *Context) GetUser(id uuid.UUID) (*domain.User, error) {
	logrus.Debugf("Application: getting user by id=%s", id)
	if user, err := c.UserRepo.GetUserByID(id); err != nil {
		logrus.Warnf("Application: cant get user with id=%s, error=%s", id, err)
		return nil, errors.New("database error")
	} else {
		if user == nil {
			return nil, nil
		}

		logrus.Debugf("Application: got user '%s' by id=%s", user.Login, id)
		return user, nil
	}
}

func (c *Context) UpdateUserLogin(id uuid.UUID, newLogin string) (bool, error) {
	logrus.Debugf("Application: updating user by id=%s login to '%s'", id, newLogin)
	// get user
	if user, getErr := c.UserRepo.GetUserByID(id); getErr == nil {
		user.Login = newLogin
		//validate
		if validateErr := user.Validate(); validateErr == nil {
			// save (update)
			if ok, saveErr := c.UserRepo.UpdateUser(*user); saveErr == nil {
				return ok, nil
			} else {
				logrus.Errorf("Application: user %s update error: %s", id, saveErr)
				// hide real db error
				return false, errors.New("database error")
			}
		} else {
			logrus.Errorf("Application: user %s validatation error: %s", id, validateErr)
			return false, validateErr
		}
	} else {
		logrus.Errorf("Application: user %s get error: %s", id, getErr)
		return false, getErr
	}
}
