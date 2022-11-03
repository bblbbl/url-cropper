package srv

import (
	"context"
	"urls/internal/repo"
	"urls/internal/service"
	cropper "urls/pkg/rpc/proto"
)

type CropperServer struct {
	cropper.UnimplementedUrlCropperServer
	urlService service.UrlService
}

func NewCropperServer() *CropperServer {
	return &CropperServer{
		urlService: service.NewUrlService(repo.NewMysqlUrlRepo()),
	}
}

func (s *CropperServer) CropUrl(_ context.Context, rq *cropper.CropRequest) (*cropper.CroppedUrl, error) {
	url := rq.GetUrl()

	result := &cropper.CroppedUrl{Url: s.urlService.CropUrl(url)}
	return result, nil
}
