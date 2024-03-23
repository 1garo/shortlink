package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName       string
	DbCollection string
	DbUri        string
	Addr         string
}

func NewConfig(filename ...string) (Config, error) {
	if len(filename) == 0 {
		filename = append(filename, ".env")
	}

	if len(filename) > 1 {
		return Config{}, fmt.Errorf("cannot pass more than 1 filename")
	}

	if err := godotenv.Load(filename...); err != nil {
		return Config{}, fmt.Errorf("no .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return Config{}, fmt.Errorf("MONGODB_URI not set")
	}

	dbName := os.Getenv("DATABASE")
	if dbName == "" {
		return Config{}, fmt.Errorf("DATABASE not set")
	}
	collectionName := os.Getenv("COLLECTION")
	if collectionName == "" {
		return Config{}, fmt.Errorf("COLLECTION not set")
	}

	addr := os.Getenv("PORT")
	if addr == "" {
		return Config{}, fmt.Errorf("PORT not set")
	}

	return Config{
		DbName:       dbName,
		DbCollection: collectionName,
		DbUri:        uri,
		Addr:         addr,
	}, nil
}
