package servers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserServerUpdateArgs struct {
	ID    uuid.UUID `json:"id"`
	Login string    `json:"login"`
}

type UserServerUpdateResult struct {
	IsOk bool `json:"is_ok"`
}

func (rpc *UserServer) UpdateUser(r *http.Request, args *UserServerUpdateArgs, result *UserServerUpdateResult) error {
	logrus.Debugf("Updating user %s with login '%s'", args.ID, args.Login)

	ok, err := rpc.AppContext.UpdateUserLogin(args.ID, args.Login)
	if err == nil {
		*result = UserServerUpdateResult{
			IsOk: ok,
		}
	}

	return err
}
