package servers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserServerGetArgs struct {
	ID uuid.UUID
}

type userView struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	RegisteredAt time.Time `json:"registered_at,string"`
}

type UserServerGetResult struct {
	User userView `json:"user,omitempty"`
}

func (rpc *UserServer) GetUser(r *http.Request, args *UserServerGetArgs, result *UserServerGetResult) error {
	logrus.Debugf("Getting user by id=%s", args.ID)

	user, err := rpc.AppContext.GetUser(args.ID)
	if err == nil {
		*result = UserServerGetResult{
			User: userView{
				ID:           user.ID,
				Login:        user.Login,
				RegisteredAt: user.RegisteredAt,
			},
		}
	}
	return err
}
