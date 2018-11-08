package main

import (
	"fmt"
	"net/http"

	"github.com/funvit/go-user-rpc-demo/infrastructure/repository"
	"github.com/funvit/go-user-rpc-demo/servers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sirupsen/logrus"
)

const port int = 1234

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	userServer := servers.NewUserServer(
		// &repository.UserRepositoryMongodb{},
		repository.NewUserRepositoryInmem(),
	)

	s.RegisterService(userServer, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	fmt.Printf("Json-RPC on localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		panic(err)
	}
}
