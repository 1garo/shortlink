package shortlink

import (
	"context"
	"errors"
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
	cfg        config.Config
	collection *mongo.Collection
}

func NewShortLinkService(client *mongo.Client, config config.Config) *ShortLinkService {
	coll := client.Database(config.DbName).Collection(config.DbCollection)
	return &ShortLinkService{
		client,
		config,
		coll,
	}
}

func (s *ShortLinkService) Redirect(url string) (TinyUrlSchema, error) {
	var result TinyUrlSchema
	filter := bson.D{{"$text", bson.D{{"$search", url}}}}
	update := bson.D{{"$inc", bson.D{{"count", 1}}}}
	err := s.collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	return result, err
}

func (s *ShortLinkService) ShortenUrlHandler(inputUrl string) (string, error) {
	if !util.IsValidUrl(inputUrl) {
		log.Println("[ShortenUrlHandler]: bad url: should have http or https")
		return "", &service.ServiceError{
			Err:  errors.New("bad url: should have http or https"),
			Code: http.StatusBadRequest,
		}
	}

	url := util.GenerateRandomShortURL(s.conn, s.cfg)

	doc := bson.D{{"shortUrl", url}, {"count", 0}, {"originalUrl", inputUrl}}
	_, err := s.collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Println("[ShortenUrlHandler]: %w", err)
		return "", &service.ServiceError{
			Err:  errors.New("InternalServerError"),
			Code: http.StatusInternalServerError,
		}
	}

	return url, &service.ServiceError{}
}
