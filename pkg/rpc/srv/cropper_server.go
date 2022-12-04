package srv

import (
	"context"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/internal/service"
	cropper "urls/pkg/rpc/proto"
)

type CropperServer struct {
	cropper.UnimplementedUrlCropperServer
	urlService *service.UrlService
}

func NewCropperServer(ctx context.Context, readRepo repo.UrlReadRepo, writeRepo repo.UrlWriteRepo, producer messaging.UrlProducer) *CropperServer {
	return &CropperServer{
		urlService: service.NewUrlService(ctx).WithUrlReadRepo(readRepo).WithUrlWriteRepo(writeRepo).WithProducer(producer),
	}
}

func (s *CropperServer) CropUrl(_ context.Context, rq *cropper.CropRequest) (*cropper.CroppedUrl, error) {
	url := rq.GetUrl()

	result := &cropper.CroppedUrl{Url: s.urlService.CropUrl(url)}
	return result, nil
}
