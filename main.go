package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os/signal"
	"syscall"

	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	numCharsShortLink = 7
	alphabet          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type GenerateUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type GenerateUrlRequest struct {
	Url string `json:"long_url"`
}

func generateRandomShortURL(client *mongo.Client) string {
	result := make([]byte, numCharsShortLink)
	for {
		for i := 0; i < numCharsShortLink; i++ {
			randomIndex := random.Intn(len(alphabet))
			result[i] = alphabet[randomIndex]
		}
		shortLink := string(result)
		// Check if the short link isn't already used
		if !checkShortLinkExists(client, shortLink) {
			return shortLink
		}
	}
}

func checkShortLinkExists(client *mongo.Client, shortUrl string) bool {
	var err error
	// TODO: move this to a config function
	dbName := os.Getenv("DATABASE")
	if dbName == "" {
		log.Fatal("DATABASE not set.")
	}
	collectionName := os.Getenv("COLLECTION")
	if collectionName == "" {
		log.Fatal("COLLECTION not set.")
	}

	coll := client.Database(dbName).Collection(collectionName)
	var result bson.M
	err = coll.FindOne(context.Background(), bson.D{{"shortUrl", shortUrl}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Printf("No document was found with the following shortUrl: %s\n", shortUrl)
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set.")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// TODO: use gin instead of mux
	r := gin.Default()

	// Define your redirect handler
	redirectHandler := func(c *gin.Context) {
		// Perform the redirection with a 302 status code
		c.Redirect(http.StatusFound, "https://google.com")
	}

	shortenUrl := func(c *gin.Context) {
		c.Header("Content-Type", "application/json")


		var input GenerateUrlRequest
		var err error

		if err = json.NewDecoder(c.Request.Body).Decode(&input); err != nil {
			log.Println("could not decode body.")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "InternalServerError",
			})
			return
		}

		// TODO: insert new collection with long_url

		// Generate randon short url
		str := generateRandomShortURL(client)

		var output GenerateUrlResponse

		output.ShortUrl = str
		c.JSON(http.StatusOK, output)
	}

	// Attach the redirect handler to a specific path
	r.GET("/old-url", redirectHandler)
	r.POST("/shorten", shortenUrl)

	// Create a server instance with the router
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
	// shortURL := generateRandomShortURL(client)
	// fmt.Println(shortURL)
}
