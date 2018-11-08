package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	apprpc "github.com/funvit/go-user-rpc-demo/rpc"
)

func main() {
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	userRPC := new(apprpc.UserRpc)

	s.RegisterService(userRPC, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":1234", r)
}
