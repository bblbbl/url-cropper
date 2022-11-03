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

type urlRequest struct {
	Url string `form:"url" json:"url" uri:"hash" binding:"required"`
}

func (uh UrlHandler) Crop(ctx *gin.Context) {
	var request urlRequest
	err := ctx.ShouldBind(&request)
	if err != nil || request.Url == "" {
		ctx.JSON(422, gin.H{
			"message": "url is required filed",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"url": uh.urlService.CropUrl(request.Url),
	})
}

func (uh UrlHandler) Redirect(ctx *gin.Context) {
	var request urlRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(422, gin.H{
			"message": "hash param is required",
		})
	}

	url, err := uh.urlService.GetLongUrl(request.Url)
	if err != nil {
		ctx.Status(404)
		return
	}

	ctx.Redirect(302, url)
}
