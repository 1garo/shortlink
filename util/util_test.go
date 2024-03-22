package util

import (
	"context"
	"testing"

	"github.com/1garo/shortlink/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func TestIsValidUrl(t *testing.T) {
	testCases := []struct {
		url  string
		want bool
	}{
		{"https://www.google.com", true},
		{"http://www.google.com", true},
		{"http:www.google.com", false},
		{"www.google.com", false},
	}

	for _, tt := range testCases {
		ok := IsValidUrl(tt.url)
		assert.Equal(t, tt.want, ok)
	}
}

func TestGenerateRandomShortURL(t *testing.T) {
	cfg, err := config.NewConfig("../.env.test")

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

	url := GenerateRandomShortURL(client, cfg)

	assert.Equal(t, len(url), 7)
}

func TestCheckShortLinkExists(t *testing.T) {
	cfg, err := config.NewConfig("../.env.test")
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
	collection := client.Database(cfg.DbName).Collection(cfg.DbCollection)

	testCases := []struct {
		url  string
		want bool
	}{
		{"testShortUrl", true},
		{"urlNotFound", false},
	}

	for _, tt := range testCases {
		ok := checkShortLinkExists(collection, tt.url)

		assert.Equal(t, tt.want, ok)
	}
}
