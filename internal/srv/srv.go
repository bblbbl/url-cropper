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

func InitServer(ctx context.Context, readRepo repo.UrlReadRepo, writeRepo repo.UrlWriteRepo, producer messaging.UrlProducer) *gin.Engine {
	server := gin.Default()
	err := server.SetTrustedProxies([]string{})
	if err != nil {
		etc.GetLogger().Fatalf("failed set trust proxies. err: %s\n", err)
	}

	server.POST("/crop", NewUrlHandlerCrop(ctx, readRepo, writeRepo, producer).Crop)
	server.GET("/go/:hash", NewUrlHandlerRedirect(ctx, readRepo).Redirect)

	return server
}

func InitRpc(ctx context.Context, readRepo repo.UrlReadRepo, writeRepo repo.UrlWriteRepo, producer messaging.UrlProducer) *grpc.Server {
	server := grpc.NewServer()
	cropperServer := srv.NewCropperServer(ctx, readRepo, writeRepo, producer)

	cropper.RegisterUrlCropperServer(server, cropperServer)

	return server
}
