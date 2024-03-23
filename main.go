package main

import (
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"log"
	"os"

	"github.com/1garo/shortlink/cmd/api"
	"github.com/1garo/shortlink/config"
	"github.com/1garo/shortlink/db"

)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}

	client := db.DbConnect(cfg.DbUri)
	defer db.DbDisconnect(client)

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
