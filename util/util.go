package util

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/1garo/shortlink/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	numCharsShortLink = 7
	alphabet          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func checkShortLinkExists(collection *mongo.Collection, shortUrl string) bool {
	filter := bson.D{{"$text", bson.D{{"$search", shortUrl}}}}
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

func GenerateRandomShortURL(client *mongo.Client, config config.Config) string {
	result := make([]byte, numCharsShortLink)
	coll := client.Database(config.DbName).Collection(config.DbCollection)
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

func IsValidUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func GracefulShutdown(srv *http.Server) {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server stopped.")
}

func SetupUrlTest(collection *mongo.Collection) error {
	filter := bson.D{{"shortUrl", "testShortUrl"}}
	update := bson.D{
		{"$set", bson.D{
			{"shortUrl", "testShortUrl"},
			{"count", 0},
			{"originalUrl", "https://www.google.com"},
		}},
	}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}
