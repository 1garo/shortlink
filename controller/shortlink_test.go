package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/service/db"
	"github.com/1garo/shortlink/service/shortlink"
	"github.com/1garo/shortlink/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestShortenUrl(t *testing.T) {
	cfg, err := config.NewConfig("../.env.test")
	assert.Nil(t, err)

	client := db.DbConnect(cfg.DbUrl)
	defer t.Cleanup(func() {
		err := db.DbCleanup(client, cfg)
		assert.Nil(t, err)
	})

	router := SetupHandler(client, cfg)

	testCases := []struct {
		body string
		code int
	}{
		{`{
			"url": "https://www.google.com"
		}`, http.StatusOK},
		{`{
			"url": "htt://www.google.com"
		}`, http.StatusBadRequest},
		{`{}`, http.StatusBadRequest},
	}

	for _, tt := range testCases {
		w := httptest.NewRecorder()
		jsonBody := []byte(tt.body)
		bodyReader := bytes.NewReader(jsonBody)
		req, _ := http.NewRequest("POST", "/shorten", bodyReader)
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestRedirectHandler(t *testing.T) {
	cfg, err := config.NewConfig("../.env.test")
	assert.Nil(t, err)

	client := db.DbConnect(cfg.DbUrl)
	defer t.Cleanup(func() {
		err := db.DbCleanup(client, cfg)
		assert.Nil(t, err)
	})
	collection := client.Database(cfg.DbName).Collection(cfg.DbCollection)
	err = util.SetupUrlTest(collection)
	assert.Nil(t, err)

	router := SetupHandler(client, cfg)

	testCases := []struct {
		uri         string
		expectedUrl string
		code        int
		count       int
	}{
		{"/testShortUrl", "https://www.google.com", http.StatusFound, 1},
		{"/not_found", "", http.StatusNotFound, 0},
	}

	for _, tt := range testCases {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", tt.uri, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
		assert.Equal(t, tt.expectedUrl, w.Result().Header.Get("Location"))

		filter := bson.D{{"$text", bson.D{{"$search", strings.TrimLeft(tt.uri, "/")}}}}
		var result shortlink.TinyUrlSchema
		_ = collection.FindOne(context.Background(), filter).Decode(&result)

		assert.Equal(t, tt.count, result.Count)
	}
}
