package shortlink

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/service"
	"github.com/1garo/shortlink/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShortLinkService struct {
	conn       *mongo.Client
	cfg        *config.Config
	collection *mongo.Collection
}

func NewShortLinkService(client *mongo.Client, config *config.Config) *ShortLinkService {
	coll := client.Database(config.DbName).Collection(config.DbCollection)
	return &ShortLinkService{
		client,
		config,
		coll,
	}
}

func (s *ShortLinkService) Redirect(url string) (result TinyUrlSchema, err error) {
	filter := bson.D{{"$text", bson.D{{"$search", url}}}}
	update := bson.D{{"$inc", bson.D{{"count", 1}}}}
	err = s.collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		errMsg := fmt.Sprintf("No document was found with the following url: %s", url)
		log.Println("[RedirectHandler]: ", errMsg)
		err = &service.ServiceError{
			Err:  errors.New(errMsg),
			Code: http.StatusNotFound,
		}
		return
	} else if err != nil {
		log.Printf("[RedirectHandler]: InternalServerError: %s\n", url)
		err = &service.ServiceError{
			Err:  errors.New("InternalServerError"),
			Code: http.StatusInternalServerError,
		}
		return
	}

	return
}

func (s *ShortLinkService) ShortenUrl(inputUrl string) (url string, err error) {
	if !util.IsValidUrl(inputUrl) {
		log.Println("[ShortenUrlHandler]: bad url: should have http or https")
		err = &service.ServiceError{
			Err:  errors.New("bad url: should have http or https"),
			Code: http.StatusBadRequest,
		}
		return
	}

	url = util.GenerateRandomShortURL(s.conn, s.collection)

	doc := bson.D{{"shortUrl", url}, {"count", 0}, {"originalUrl", inputUrl}}
	_, err = s.collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("[ShortenUrlHandler]: %s\n", err)
		err = &service.ServiceError{
			Err:  errors.New("InternalServerError"),
			Code: http.StatusInternalServerError,
		}
	}

	return
}
