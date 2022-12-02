package srv

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"urls/internal/service"
	"urls/pkg/etc"
	cropper "urls/pkg/rpc/proto"
	"urls/pkg/rpc/srv"
)

func InitServer(we *service.WriteExecutor, ctx context.Context) *gin.Engine {
	server := gin.Default()
	err := server.SetTrustedProxies([]string{})
	if err != nil {
		etc.GetLogger().Fatalf("failed set trust proxies. err: %s\n", err)
	}

	server.POST("/crop", NewUrlHandler(we, ctx).Crop)
	server.GET("/go/:hash", NewUrlHandler(we, ctx).Redirect)

	return server
}

func InitRpc(we *service.WriteExecutor, ctx context.Context) *grpc.Server {
	server := grpc.NewServer()
	cropperServer := srv.NewCropperServer(we, ctx)

	cropper.RegisterUrlCropperServer(server, cropperServer)

	return server
}
