package util

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/1garo/shortlink/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	numCharsShortLink = 7
	alphabet          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func checkShortLinkExists(collection *mongo.Collection, shortUrl string) bool {
	filter := bson.D{{"shortUrl", shortUrl}}
	var _result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&_result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("No document was found with the following shortUrl: %s\n", shortUrl)
		return false
	} else if err != nil {
		panic(err)
	}

	return true
}

func GenerateRandomShortURL(client *mongo.Client) string {
	result := make([]byte, numCharsShortLink)
	// TODO: move this to a config function

	coll := client.Database(config.Cfg.DbName).Collection(config.Cfg.DbCollection)
	for {
		for i := 0; i < numCharsShortLink; i++ {
			randomIndex := random.Intn(len(alphabet))
			result[i] = alphabet[randomIndex]
		}
		shortLink := string(result)
		// Check if the short link isn't already used
		if !checkShortLinkExists(coll, shortLink) {
			return shortLink
		}
	}
}
