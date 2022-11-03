package service

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"urls/internal/repo"
	"urls/pkg/etc"
)

func TestUrlService_CropUrl(t *testing.T) {
	inti()
	service := NewUrlService(repo.NewMysqlUrlRepo())

	cropped := service.CropUrl("https://google.com")

	assert.Equal(t, "http://localhost:8000/go/2Ep", cropped)
}

func TestGetLongUrl(t *testing.T) {
	inti()
	service := NewUrlService(repo.NewMysqlUrlRepo())

	url := strconv.Itoa(rand.Intn(100000))

	cropped := service.CropUrl(url)
	urlParts := strings.Split(cropped, "/")
	hash := urlParts[len(urlParts)-1]

	longUrl, err := service.GetLongUrl(hash)
	if err != nil {
		log.Fatal("failed to get long url by hash")
	}

	assert.Equal(t, longUrl, url)
}

func TestUrlService_buildFullShortUrl(t *testing.T) {
	inti()
	service := NewUrlService(repo.NewMysqlUrlRepo())

	hash := "12345"
	url := service.buildFullShortUrl(hash)

	cnf := etc.GetConfig()

	assert.Equal(t, url, fmt.Sprintf("%s://%s/go/%s", cnf.Http.Schema, cnf.App.Host, hash))
}

func inti() {
	path, err := filepath.Abs("../../.env.test")
	if err != nil {
		log.Fatal("failed to get root path")
	}

	err = godotenv.Load(path)
	if err != nil {
		log.Fatal("failed to load .env")
	}

	etc.InitLogger()
	etc.InitConfig()
}
