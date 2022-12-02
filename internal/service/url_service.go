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

var (
	generator     *HashGenerator
	generatorOnce sync.Once
)

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

func Generator() *HashGenerator {
	generatorOnce.Do(func() {
		generator = &HashGenerator{}
	})

	return generator
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

func (us *UrlService) CropUrl(url string) string {
	if v, ok := us.cache.HashByUrl(url); ok {
		return us.buildFullShortUrl(v)
	}

	existUrl := us.urlRepo.GetByFull(url)
	if existUrl != nil {
		return existUrl.Hash
	}

	hash := us.createUrlHash()
	us.cache.PutUrl(hash, url)

	us.executor.JobChan <- CreateUrlJob{
		url,
		hash,
	}

	return us.buildFullShortUrl(hash)
}

func (us *UrlService) GetLongUrl(hash string) (string, error) {
	//if v, ok := us.cache.GetUrl(hash); ok {
	//	return v, nil
	//}

	url := us.urlRepo.GetByHash(hash)
	if url == nil {
		return "", errors.New("short url not found")
	}

	return url.Long, nil
}

func (us *UrlService) createUrlHash() string {
	g := Generator()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.lastId == 0 {
		g.lastId = us.urlRepo.GetLastId()
	}

	hash, _ := us.data.Encode([]int{g.lastId + 1})
	g.lastId += 1

	return hash
}

func (us *UrlService) buildFullShortUrl(hash string) string {
	return fmt.Sprintf("%s://%s/go/%s", us.cnf.Http.Schema, us.cnf.App.Host, hash)
}
