package srv

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"urls/internal/repo"
	"urls/internal/service"
)

type UrlHandler struct {
	urlService service.UrlService
}

func NewUrlHandler(we *service.WriteExecutor) UrlHandler {
	return UrlHandler{
		urlService: service.NewUrlService(repo.NewMysqlUrlRepo(), we),
	}
}

type urlRequest struct {
	Url string `form:"url" json:"url" uri:"hash" binding:"required"`
}

func (uh UrlHandler) Crop(ctx *gin.Context) {
	var request urlRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil || request.Url == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "url is required filed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": uh.urlService.CropUrl(request.Url),
	})
}

func (uh UrlHandler) Redirect(ctx *gin.Context) {
	var request urlRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "hash param is required",
		})
	}

	url, err := uh.urlService.GetLongUrl(request.Url)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusFound, url)
}
