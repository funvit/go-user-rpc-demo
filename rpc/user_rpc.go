package rpc

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/funvit/go-user-rpc-demo/application"
)

type UserRpc struct {
	AppContext application.Context
}

type UserRpcAddArgs struct {
	Login string
}

type UserRpcAddResult struct {
	NewUserID uuid.UUID
	Err       string
}

func (rpc *UserRpc) AddUser(r *http.Request, args *UserRpcAddArgs, result *UserRpcAddResult) error {
	log.Printf("Adding user with login %s\n", args.Login)

	user, err := rpc.AppContext.AddUser(args.Login)
	if err != nil {
		return err
	}

	*result = UserRpcAddResult{
		NewUserID: user.ID,
		Err:       err.Error(),
	}

	return nil
}
