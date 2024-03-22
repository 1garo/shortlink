package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1garo/shortlink/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestShortenUrl(t *testing.T) {
	cfg, err := config.NewConfig("../../.env.test")

	assert.Nil(t, err)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DbUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

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
	router := SetupRouter(client, cfg)

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
	cfg, err := config.NewConfig("../../.env.test")

	assert.Nil(t, err)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DbUri))
	if err != nil {
		panic(err)
	}
	// TODO: used in some places, make a function
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	testCases := []struct {
		uri string
		expectedUrl string
		code int
	}{
		{"/testShortUrl", "https://www.google.com", http.StatusFound},
		{"/badurl", "", http.StatusBadRequest},
	}
	router := SetupRouter(client, cfg)

	for _, tt := range testCases {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", tt.uri, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
		assert.Equal(t, tt.expectedUrl, w.Result().Header.Get("Location"))

	}
}
