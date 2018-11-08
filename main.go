package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/funvit/go-user-rpc-demo/infrastructure/repository"
	"github.com/funvit/go-user-rpc-demo/servers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/sirupsen/logrus"
)

func main() {
	// flags
	var debugMode = flag.Bool("d", false, "debug mode")
	var useInmemDB = flag.Bool("in-mem", false, "use im-mem db")
	var mongodbConnStr = flag.String("mongo-conn", "", "mongodb connection string")
	var mongodbDatabaseName = flag.String("mongo-database-name", "user_rpc", "mongodb database name")
	var rpcPort = flag.Int("p", 1234, "RPC port")

	flag.Parse()

	// check flags
	if (*useInmemDB && *mongodbConnStr != "") || (!*useInmemDB && *mongodbConnStr == "") {
		fmt.Println("Run with -in-mem or with -mongo-conn!")
		os.Exit(1)
	}

	// main
	if *debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	var userServer *servers.UserServer
	if *useInmemDB {
		logrus.Info("Using in-mem as DB")
		userServer = servers.NewUserServer(
			repository.NewUserRepositoryInmem(),
		)
	}
	if *mongodbConnStr != "" {
		logrus.Infof("Using mongodb %s with database '%s'", *mongodbConnStr, *mongodbDatabaseName)
		userServer = servers.NewUserServer(
			repository.NewUserRepositoryMongo(
				*mongodbConnStr,
				*mongodbDatabaseName,
				time.Second*2,
			),
		)
	}

	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	s.RegisterService(userServer, "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)
	logrus.Infof("Json-RPC on localhost:%d/rpc\n", *rpcPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *rpcPort), r); err != nil {
		panic(err)
	}
}
