package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/1garo/shortlink/config"
	util "github.com/1garo/shortlink/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	client *mongo.Client
}

func NewHandler(client *mongo.Client) *Handler {
	return &Handler{
		client,
	}
}

func SetupRouter(client *mongo.Client) *gin.Engine {
	r := gin.Default()
	h := NewHandler(client)
	r.GET("/:url", h.RedirectHandler)
	r.POST("/shorten", h.ShortenUrl)

	return r
}

func (h *Handler) RedirectHandler(c *gin.Context) {
	url := c.Param("url")
	coll := h.client.Database(config.Cfg.DbName).Collection(config.Cfg.DbCollection)

	log.Println(url)

	var result TinyUrlSchema
	filter := bson.D{{"shortUrl", url}}
	update := bson.D{{"$inc", bson.D{{"count", 1}}}}
	err := coll.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if err == mongo.ErrNoDocuments {
		errMsg := fmt.Sprintf("No document was found with the following url: %s", url)
		log.Println(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		return
	}
	if err != nil {
		errMsg := fmt.Sprintf("InternalServerError: %s", url)
		log.Println(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "InternalServerError",
		})
		return
	}

	c.Redirect(http.StatusFound, result.OriginalUrl)
}

func (h *Handler) ShortenUrl(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	coll := h.client.Database(config.Cfg.DbName).Collection(config.Cfg.DbCollection)

	var input ShortenUrlRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("could not decode body.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !strings.HasPrefix(input.Url, "http://") && !strings.HasPrefix(input.Url, "https://") {
		log.Println("bad prefix")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad prefix: should have http or https",
		})
		return

	}

	url := util.GenerateRandomShortURL(h.client)

	doc := bson.D{{"shortUrl", url}, {"count", 0}, {"originalUrl", input.Url}}
	_, err := coll.InsertOne(context.Background(), doc)
	if err != nil {
		log.Println("could not insert new document: %w", err)
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
