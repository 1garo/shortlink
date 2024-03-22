package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"log"
	"os"

	"github.com/1garo/shortlink/cmd/api"
	"github.com/1garo/shortlink/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DbUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Addr),
		Handler: api.SetupRouter(client, cfg),
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
	// block until signal is received
	<-quit
	//GracefulShutdown(srv)
}
