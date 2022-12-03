package srv

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/etc"
	cropper "urls/pkg/rpc/proto"
	"urls/pkg/rpc/srv"
)

func InitServer(ctx context.Context, repo repo.UrlRepo, producer messaging.UrlProducer) *gin.Engine {
	server := gin.Default()
	err := server.SetTrustedProxies([]string{})
	if err != nil {
		etc.GetLogger().Fatalf("failed set trust proxies. err: %s\n", err)
	}

	server.POST("/crop", NewUrlHandlerCrop(ctx, repo, producer).Crop)
	server.GET("/go/:hash", NewUrlHandlerRedirect(ctx, repo).Redirect)

	return server
}

func InitRpc(ctx context.Context, repo repo.UrlRepo, producer messaging.UrlProducer) *grpc.Server {
	server := grpc.NewServer()
	cropperServer := srv.NewCropperServer(ctx, repo, producer)

	cropper.RegisterUrlCropperServer(server, cropperServer)

	return server
}
