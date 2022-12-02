package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"os/signal"
	"syscall"
	"urls/internal/service"
	"urls/internal/srv"
	"urls/pkg/database"
	"urls/pkg/etc"
)

const realiseMode = "release"

func main() {
	etc.InitLogger()

	cnf := etc.GetConfig()
	if cnf.App.Mode == realiseMode {
		gin.SetMode(gin.ReleaseMode)
		etc.GetLogger().Info("application run in realise mode")
	}

	database.GetConnection()

	ctx, cancelFunc := context.WithCancel(context.Background())
	writeExecutor := service.NewWriteExecutor(ctx).Start()

	go func() {
		rpcServer := srv.InitRpc(writeExecutor, ctx)
		l, err := net.Listen(cnf.Rpc.Network, fmt.Sprintf(":%d", cnf.Rpc.Port))
		if err != nil {
			panic(err)
		}

		if err = rpcServer.Serve(l); err != nil {
			panic(err)
		}
	}()

	go func() {
		server := srv.InitServer(writeExecutor, ctx)
		if err := server.Run(fmt.Sprintf(":%d", cnf.Http.Port)); err != nil {
			panic(err)
		}
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-term
	terminate(cancelFunc)
}

func terminate(cancelFunc context.CancelFunc) {
	etc.GetLogger().Info("start shutting down server")

	database.CloseRedisConnection()
	database.CloseMysqlConnection()
	etc.FlushLogger()

	cancelFunc()

	etc.GetLogger().Info("server successful shutting down")
}
