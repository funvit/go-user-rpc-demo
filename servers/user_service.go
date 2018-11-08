package servers

import (
	"github.com/funvit/go-user-rpc-demo/application"
	"github.com/funvit/go-user-rpc-demo/infrastructure/repository"
)

type UserServer struct {
	AppContext application.Context
}

func NewUserServer(userRepo repository.UserRepository) *UserServer {
	return &UserServer{
		AppContext: application.Context{
			UserRepo: userRepo,
		},
	}
}
