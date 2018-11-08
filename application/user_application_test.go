package application_test

import (
	"testing"

	"github.com/funvit/go-user-rpc-demo/application"
	"github.com/funvit/go-user-rpc-demo/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.PanicLevel)
}

func TestUserApplicationAddGetUpdate(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	appCtx := application.Context{
		UserRepo: repository.NewUserRepositoryInmem(),
	}
	var _userId uuid.UUID

	// add
	user, err := appCtx.AddUser("test_login_1")
	if err != nil {
		t.Error(err)
	} else {
		if user == nil {
			t.Error("returned user as nil")
		} else {
			_userId = user.ID
		}
	}

	// get
	if storedUser, getErr := appCtx.GetUser(_userId); getErr != nil {
		t.Errorf("cant get stored user by id=%s, error=%s", _userId, getErr)
	} else {
		if storedUser == nil {
			t.Error("got user as nil")
		}
	}

	// update
	newLogin := "test_user_l_alt"
	if ok, updateErr := appCtx.UpdateUserLogin(_userId, newLogin); updateErr != nil {
		t.Error("user update error (Login)")
	} else {
		if !ok {
			t.Error("user update returned false")
		} else {
			// re-check by get
			if userAfterUpdate, getErr := appCtx.GetUser(_userId); getErr == nil {
				if userAfterUpdate.Login != newLogin {
					t.Errorf("user get after update Login returned wrong value '%s' (expected: %s)", userAfterUpdate.Login, newLogin)
				}
			} else {
				t.Error("cant re-read user")
			}
		}
	}
}

func TestUserApplicationGetNotExist(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	appCtx := application.Context{
		UserRepo: repository.NewUserRepositoryInmem(),
	}
	userId := uuid.New()

	user, err := appCtx.GetUser(userId)
	if err == nil {
		t.Error("get non-exist user must return error")
	}
	if user != nil {
		t.Error("get not-exist user must return nil as User")
	}
}

func TestUserApplicationAddWithEmptyLoginMustReturnError(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	appCtx := application.Context{
		UserRepo: repository.NewUserRepositoryInmem(),
	}

	// add
	if user, addErr := appCtx.AddUser(""); addErr == nil {
		t.Error("must return error")
	} else {
		if user != nil {
			t.Error("must return user as nil")
		}
	}
}
