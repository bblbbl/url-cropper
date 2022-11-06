package srv

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"urls/internal/service"
	"urls/pkg/etc"
	cropper "urls/pkg/rpc/proto"
	"urls/pkg/rpc/srv"
)

func InitServer(we *service.WriteExecutor) *gin.Engine {
	server := gin.Default()
	err := server.SetTrustedProxies([]string{})
	if err != nil {
		etc.GetLogger().Fatalf("failed set trust proxies. err: %s\n", err)
	}

	server.POST("/crop", NewUrlHandler(we).Crop)
	server.GET("/go/:hash", NewUrlHandler(we).Redirect)

	return server
}

func InitRpc(we *service.WriteExecutor) *grpc.Server {
	server := grpc.NewServer()
	cropperServer := srv.NewCropperServer(we)

	cropper.RegisterUrlCropperServer(server, cropperServer)

	return server
}
