package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/speps/go-hashids/v2"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/etc"
	"urls/pkg/utils"
)

type UrlService struct {
	producer     messaging.UrlProducer
	data         *hashids.HashID
	cache        repo.UrlCacheRepo
	urlReadRepo  repo.UrlReadRepo
	urlWriteRepo repo.UrlWriteRepo
	cnf          *etc.Config
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

	existUrl := us.urlReadRepo.GetByFull(url)
	if existUrl != nil {
		return existUrl.Hash
	}

	hash := generateShortLink(url)
	us.cache.PutUrl(hash, url)

	urlModel := repo.NewUrl(hash, url)
	err := us.producer.PutUrlMessage(urlModel)
	if err != nil {
		_ = us.urlWriteRepo.CreateUrl(urlModel)
	}

	return us.buildFullShortUrl(hash)
}

func (us *UrlService) GetLongUrl(hash string) (string, error) {
	if v, ok := us.cache.GetUrl(hash); ok {
		return v, nil
	}

	url := us.urlReadRepo.GetByHash(hash)
	if url == nil {
		return "", errors.New("short url not found")
	}

	return url.Long, nil
}

func (us *UrlService) buildFullShortUrl(hash string) string {
	return fmt.Sprintf("%s://%s/go/%s", us.cnf.Http.Schema, us.cnf.App.Host, hash)
}

func (us *UrlService) WithUrlReadRepo(urlRepo repo.UrlReadRepo) *UrlService {
	us.urlReadRepo = urlRepo

	return us
}

func (us *UrlService) WithUrlWriteRepo(urlRepo repo.UrlWriteRepo) *UrlService {
	us.urlWriteRepo = urlRepo

	return us
}

func (us *UrlService) WithProducer(producer messaging.UrlProducer) *UrlService {
	us.producer = producer

	return us
}

func generateShortLink(initialLink string) string {
	b := sha256.Sum256([]byte(initialLink))
	return utils.B2S(b[:8])
}
