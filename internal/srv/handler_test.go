package srv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"testing"
	"urls/pkg/database"
	"urls/pkg/etc"
)

func TestUrlHandler_Crop(t *testing.T) {
	server := initServer()

	request := urlRequest{
		Url: "test",
	}

	m, _ := json.Marshal(&request)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/crop", bytes.NewBuffer(m))
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("failed to parse response")
	}

	if val, ok := res["url"]; !ok {
		t.Fatal("response missing url filed")
	} else {
		cnf := etc.GetConfig()
		reg := regexp.MustCompile(fmt.Sprintf("%s:\\/\\/%s\\/go\\/(.)+", cnf.Http.Schema, cnf.App.Host))
		assert.MatchRegex(t, val.(string), reg)
	}
}

func TestUrlHandler_Crop_422(t *testing.T) {
	server := initServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/crop", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUrlHandler_Redirect(t *testing.T) {
	server := initServer()

	path := "abc"
	cnf := etc.GetConfig()
	short := fmt.Sprintf("%s://%s/go/%s", cnf.Http.Schema, cnf.App.Host, path)
	_, _ = database.GetConnection().Exec("INSERT INTO urls (`long`, short) VALUES ('test', ?)", short)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/go/%s", path), nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
}

func initServer() *gin.Engine {
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

	return InitServer()
}
