package util

import (
	"testing"

	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/db"
	"github.com/stretchr/testify/assert"
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

	client := db.DbConnect(cfg.DbUri)
	defer db.DbDisconnect(client)

	url := GenerateRandomShortURL(client, cfg)

	assert.Equal(t, len(url), 7)
}

func TestCheckShortLinkExists(t *testing.T) {
	cfg, err := config.NewConfig("../.env.test")
	assert.Nil(t, err)

	client := db.DbConnect(cfg.DbUri)
	defer t.Cleanup(func() {
		err := db.DbCleanup(client, cfg)
		assert.Nil(t, err)
	})

	collection := client.Database(cfg.DbName).Collection(cfg.DbCollection)
	err = SetupUrlTest(collection)
	assert.Nil(t, err)

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
