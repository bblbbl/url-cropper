package srv

import (
	"github.com/gin-gonic/gin"
	"urls/internal/repo"
	"urls/internal/service"
)

type UrlHandler struct {
	urlService service.UrlService
}

func NewUrlHandler() UrlHandler {
	return UrlHandler{
		urlService: service.NewUrlService(repo.NewMysqlUrlRepo()),
	}
}

func (uh UrlHandler) Crop(ctx *gin.Context) {
	url := ctx.PostForm("url")

	ctx.JSON(200, gin.H{
		"url": uh.urlService.CropUrl(url),
	})
}

func (uh UrlHandler) Redirect(ctx *gin.Context) {
	hash := ctx.Param("hash")

	url, err := uh.urlService.GetLongUrl(hash)
	if err != nil {
		ctx.Status(404)
		return
	}

	ctx.Redirect(302, url)
}
