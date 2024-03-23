package db

import (
	"context"

	"github.com/1garo/shortlink/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbConnect(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func DbDisconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func DbCleanup(client *mongo.Client, cfg config.Config) error {
	defer DbDisconnect(client)
	collection := client.Database(cfg.DbName).Collection(cfg.DbCollection)
	filter := bson.D{{}}
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}
