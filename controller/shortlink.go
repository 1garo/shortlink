package controller

import (
	"log"
	"net/http"

	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/service"
	"github.com/1garo/shortlink/service/shortlink"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	client *mongo.Client
	config *config.Config
}

func NewHandler(client *mongo.Client, config *config.Config) *Handler {
	return &Handler{
		client,
		config,
	}
}

func SetupHandler(client *mongo.Client, config *config.Config) *gin.Engine {
	r := gin.Default()
	h := NewHandler(client, config)
	r.GET("/:url", h.RedirectHandler)
	r.POST("/shorten", h.ShortenUrlHandler)

	return r
}

func (h *Handler) RedirectHandler(c *gin.Context) {
	log.Println("[RedirectHandler]")
	url := c.Param("url")

	svc := shortlink.NewShortLinkService(h.client, h.config)
	result, err := svc.Redirect(url)

	if err != nil {
		err := err.(*service.ServiceError)
		log.Printf("[ShortenUrl]: %s\n", err.Err.Error())
		c.JSON(err.Code, gin.H{
			"error": err.Err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, result.OriginalUrl)
}

func (h *Handler) ShortenUrlHandler(c *gin.Context) {
	log.Println("[ShortenUrlHandler]")
	c.Header("Content-Type", "application/json")

	var input shortlink.ShortenUrlRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("[ShortenUrl]: could not decode body.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not decode body",
		})
		return
	}

	svc := shortlink.NewShortLinkService(h.client, h.config)
	url, err := svc.ShortenUrl(input.Url)
	if err != nil {
		err := err.(*service.ServiceError)
		log.Printf("[ShortenUrl]: %s\n", err.Err.Error())
		c.JSON(err.Code, gin.H{
			"error": err.Err.Error(),
		})
		return
	}

	output := shortlink.ShortenUrlResponse{
		ShortUrl: url,
	}
	c.JSON(http.StatusOK, output)
}
