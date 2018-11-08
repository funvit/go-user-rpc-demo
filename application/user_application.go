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
	logrus.Debugf("adding user \"%s\"", login)
	user := domain.CreateNewUser(login)
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := c.UserRepo.SaveUser(*user); err != nil {
		return nil, err
	} else {
		logrus.Debugf("added user \"%s\", new ID=%s", user.Login, user.ID)
		return user, nil
	}
}

func (c *Context) GetUser(id uuid.UUID) (*domain.User, error) {
	logrus.Debugf("getting user by id=%s", id)
	if user, err := c.UserRepo.GetUserByID(id); err != nil {
		logrus.Warnf("cant get user with id=%s, error=%s", id, err)
		return nil, errors.New("user not found")
	} else {
		logrus.Debugf("got user \"%s\" by id=%s", user.Login, id)
		return user, nil
	}
}

func (c *Context) UpdateUserLogin(id uuid.UUID, newLogin string) (bool, error) {
	logrus.Debugf("updating user by id=%s login to \"%s\"", id, newLogin)
	// get user
	if user, getErr := c.UserRepo.GetUserByID(id); getErr == nil {
		user.Login = newLogin
		//validate
		if validateErr := user.Validate(); getErr == nil {
			// save (update)
			if saveErr := c.UserRepo.SaveUser(*user); saveErr == nil {
				return true, nil
			} else {
				logrus.Errorf("user update error: %s", saveErr)
				// hide db error
				return false, errors.New("update failed")
			}
		} else {
			logrus.Errorf("user validatation error: %s", validateErr)
			return false, validateErr
		}
	} else {
		logrus.Errorf("user get error: %s", getErr)
		return false, getErr
	}
}
