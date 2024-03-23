package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	client *mongo.Client
	config config.Config
}

func NewHandler(client *mongo.Client, config config.Config) *Handler {
	return &Handler{
		client,
		config,
	}
}

func SetupRouter(client *mongo.Client, config config.Config) *gin.Engine {
	r := gin.Default()
	h := NewHandler(client, config)
	r.GET("/:url", h.RedirectHandler)
	r.POST("/shorten", h.ShortenUrl)

	return r
}

func (h *Handler) RedirectHandler(c *gin.Context) {
	// TODO: use a service here
	log.Println("[RedirectHandler]")
	url := c.Param("url")
	coll := h.client.Database(h.config.DbName).Collection(h.config.DbCollection)

	var result TinyUrlSchema
	filter := bson.D{{"$text", bson.D{{"$search", url}}}}
	update := bson.D{{"$inc", bson.D{{"count", 1}}}}
	err := coll.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		errMsg := fmt.Sprintf("No document was found with the following url: %s", url)
		log.Println("[RedirectHandler]: ", errMsg)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errMsg,
		})
		return
	} else if err != nil {
		log.Printf("[RedirectHandler]: InternalServerError: %s\n", url)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "InternalServerError",
		})
		return
	}

	c.Redirect(http.StatusFound, result.OriginalUrl)
}

func (h *Handler) ShortenUrl(c *gin.Context) {
	// TODO: use a service here
	log.Println("[ShortenUrl]")
	c.Header("Content-Type", "application/json")
	coll := h.client.Database(h.config.DbName).Collection(h.config.DbCollection)

	var input ShortenUrlRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("[ShortenUrl]: could not decode body.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not decode body",
		})
		return
	}

	if !util.IsValidUrl(input.Url) {
		log.Println("[ShortenUrl]: bad url: should have http or https")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad url: should have http or https",
		})
		return
	}

	url := util.GenerateRandomShortURL(h.client, h.config)

	doc := bson.D{{"shortUrl", url}, {"count", 0}, {"originalUrl", input.Url}}
	_, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		log.Println("[ShortenUrl]: could not insert new document: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "InternalServerError",
		})
		return
	}
	output := ShortenUrlResponse{
		ShortUrl: url,
	}
	c.JSON(http.StatusOK, output)
}
