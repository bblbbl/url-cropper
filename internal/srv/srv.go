package srv

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	cropper "urls/pkg/rpc/proto"
	"urls/pkg/rpc/srv"
)

func InitServer() *gin.Engine {
	server := gin.Default()
	err := server.SetTrustedProxies([]string{})
	if err != nil {
		log.Fatalf("failed set trust proxies. err: %s\n", err)
	}

	server.POST("/crop", NewUrlHandler().Crop)
	server.GET("/go/:hash", NewUrlHandler().Redirect)

	return server
}

func InitRpc() *grpc.Server {
	server := grpc.NewServer()
	cropperServer := srv.NewCropperServer()

	cropper.RegisterUrlCropperServer(server, cropperServer)

	return server
}
