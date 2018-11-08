package servers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type UserServerAddArgs struct {
	Login string
}

type UserServerAddResult struct {
	User userView `json:"user,omitempty"`
}

func (rpc *UserServer) AddUser(r *http.Request, args *UserServerAddArgs, result *UserServerAddResult) error {
	logrus.Debugf("UserServer: Adding user with login '%s'", args.Login)

	user, err := rpc.AppContext.AddUser(args.Login)
	if err == nil {
		*result = UserServerAddResult{
			User: userView{
				ID:           user.ID,
				Login:        user.Login,
				RegisteredAt: user.RegisteredAt,
			},
		}
	}

	return err
}
