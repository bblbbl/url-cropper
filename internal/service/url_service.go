package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"sync"
	"urls/internal/repo"
	"urls/pkg/etc"
)

var hashGenerator *HashGenerator

type UrlService struct {
	executor *WriteExecutor
	data     *hashids.HashID
	cache    repo.UrlCacheRepo
	urlRepo  repo.UrlRepo
	cnf      *etc.Config
}

type HashGenerator struct {
	lastId int
	mu     sync.RWMutex
}

func NewUrlService(urlRepo repo.UrlRepo, executor *WriteExecutor, ctx context.Context) UrlService {
	data := hashids.NewData()
	data.Salt = etc.GetConfig().Hash.Salt
	data.MinLength = 3
	hashData, _ := hashids.NewWithData(data)

	return UrlService{
		executor: executor,
		data:     hashData,
		cache:    repo.NewUrlRedisCache(ctx),
		urlRepo:  urlRepo,
		cnf:      etc.GetConfig(),
	}
}

func GetHashGenerator() *HashGenerator {
	if hashGenerator == nil {
		hashGenerator = &HashGenerator{}
	}

	return hashGenerator
}

func (us *UrlService) CropUrl(url string) string {
	if v, ok := us.cache.GetShortUrl(url); ok {
		return v
	}

	existUrl := us.urlRepo.GetByFull(url)
	if existUrl != nil {
		return existUrl.GetShort()
	}

	shortUrl := us.buildFullShortUrl(us.createUrlHash())
	us.cache.PutUrl(url, shortUrl)

	us.executor.JobChan <- CreateUrlJob{
		shortUrl, url,
	}

	return shortUrl
}

func (us *UrlService) GetLongUrl(hash string) (string, error) {
	shortUrl := us.buildFullShortUrl(hash)
	if v, ok := us.cache.GetUrl(shortUrl); ok {
		return v, nil
	}

	url := us.urlRepo.GetByShort(shortUrl)
	if url == nil {
		return "", errors.New("short url not found")
	}

	return url.GetLong(), nil
}

func (us *UrlService) createUrlHash() string {
	generator := GetHashGenerator()

	generator.mu.Lock()
	defer generator.mu.Unlock()

	if generator.lastId == 0 {
		generator.lastId = us.urlRepo.GetLastId()
	}

	hash, _ := us.data.Encode([]int{generator.lastId + 1})
	generator.lastId += 1

	return hash
}

func (us *UrlService) buildFullShortUrl(hash string) string {
	return fmt.Sprintf("%s://%s/go/%s", us.cnf.Http.Schema, us.cnf.App.Host, hash)
}
