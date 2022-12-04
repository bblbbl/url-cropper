package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/itchyny/base58-go"
	"github.com/speps/go-hashids/v2"
	"math/big"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/etc"
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

	hash, err := us.createUrlHash(url)
	if err != nil {
		panic(err)
	}

	us.cache.PutUrl(hash, url)

	urlModel := repo.NewUrl(hash, url)
	err = us.producer.PutUrlMessage(urlModel)
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

func (us *UrlService) createUrlHash(url string) (string, error) {
	hash, err := us.generateShortLink(url)
	if err != nil {
		return "", err
	}

	return hash, nil
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

func (us *UrlService) generateShortLink(initialLink string) (string, error) {
	urlHashBytes := sha256Of(initialLink)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}

	return finalString[:8], nil
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
