package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net"
	"os"
	"os/signal"
	"syscall"
	"urls/internal/srv"
	"urls/pkg/database"
	"urls/pkg/etc"
)

const realiseMode = "release"

func main() {
	etc.InitLogger()

	err := godotenv.Load()
	if err != nil {
		etc.GetLogger().Fatalf(".env read failed: %e\n", err)
	}

	etc.InitConfig()
	cnf := etc.GetConfig()

	if cnf.App.Mode == realiseMode {
		gin.SetMode(gin.ReleaseMode)
		etc.GetLogger().Info("application run in realise mode")
	}

	database.InitConnection()

	go func() {
		rpcServer := srv.InitRpc()
		l, err := net.Listen(cnf.Rpc.Network, fmt.Sprintf(":%s", cnf.Rpc.Port))
		if err != nil {
			panic(err)
		}

		if err := rpcServer.Serve(l); err != nil {
			panic(err)
		}
	}()

	go func() {
		server := srv.InitServer()
		if err = server.Run(fmt.Sprintf(":%s", cnf.Http.Port)); err != nil {
			panic(err)
		}
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-term
	terminate()
}

func terminate() {
	etc.GetLogger().Info("start shutting down server")

	database.CloseRedisConnection()
	database.CloseMysqlConnection()
	etc.FlushLogger()

	etc.GetLogger().Info("server successful shutting down")
}
