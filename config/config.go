package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DbName       string
	DbCollection string
	DbUri        string
	Addr string
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set.")
	}

	dbName := os.Getenv("DATABASE")
	if dbName == "" {
		log.Fatal("DATABASE not set.")
	}
	collectionName := os.Getenv("COLLECTION")
	if collectionName == "" {
		log.Fatal("COLLECTION not set.")
	}

	addr := os.Getenv("PORT")
	if addr == "" {
		log.Fatal("PORT not set.")
	}
	return Config{
		DbName:       dbName,
		DbCollection: collectionName,
		DbUri:        uri,
		Addr: addr,
	}
}
