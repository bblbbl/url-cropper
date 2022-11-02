package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net"
	"urls/internal/srv"
	"urls/pkg/config"
	"urls/pkg/database"
)

const realiseMode = "release"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env read failed: %e\n", err)
	}

	config.InitConfig()
	cnf := config.GetConfig()

	if cnf.App.Mode == realiseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	database.InitConnection()

	rpcServer := srv.InitRpc()
	l, err := net.Listen(cnf.Rpc.Network, fmt.Sprintf(":%s", cnf.Rpc.Port))
	if err != nil {
		panic(err)
	}

	go func(l net.Listener) {
		if err := rpcServer.Serve(l); err != nil {
			panic(err)
		}
	}(l)

	server := srv.InitServer()

	if err = server.Run(fmt.Sprintf(":%s", cnf.Http.Port)); err != nil {
		panic(err)
	}
}
