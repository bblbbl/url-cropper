package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"sync"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/etc"
)

var (
	generator     *HashGenerator
	generatorOnce sync.Once
)

type UrlService struct {
	producer messaging.UrlProducer
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

func NewUrlService(ctx context.Context) *UrlService {
	data := hashids.NewData()
	data.Salt = etc.GetConfig().Hash.Salt
	data.MinLength = 3
	hashData, _ := hashids.NewWithData(data)

	return &UrlService{
		data:  hashData,
		cache: repo.NewUrlRedisCache(ctx),
		cnf:   etc.GetConfig(),
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

	urlModel := repo.NewUrl(hash, url)
	err := us.producer.PutUrlMessage(urlModel)
	if err != nil {
		_ = us.urlRepo.CreateUrl(urlModel)
	}

	return us.buildFullShortUrl(hash)
}

func (us *UrlService) GetLongUrl(hash string) (string, error) {
	if v, ok := us.cache.GetUrl(hash); ok {
		return v, nil
	}

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

func (us *UrlService) WithUrlRepo(urlRepo repo.UrlRepo) *UrlService {
	us.urlRepo = urlRepo

	return us
}

func (us *UrlService) WithProducer(producer messaging.UrlProducer) *UrlService {
	us.producer = producer

	return us
}
