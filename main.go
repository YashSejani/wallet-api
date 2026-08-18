package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"wallet-api/api"
	"wallet-api/db/sqlc"
	"wallet-api/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}
	log.Println("successfully connected to PostgreSQL database")

	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	httpServer := &http.Server{
		Addr:    config.ServerAddress,
		Handler: server.Router(),
	}

	go func() {
		log.Printf("HTTP server running on %s", config.ServerAddress)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down HTTP server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("server exited cleanly")
}
