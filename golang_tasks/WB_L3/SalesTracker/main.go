package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"salestracker/internal/server"
	"salestracker/internal/tracker"
)

var PORT = ":8080"

func main() {
	commenterSvc := tracker.NewService()
	defer commenterSvc.Close()

	srv := &http.Server{
		Addr:    PORT,
		Handler: server.NewServer(commenterSvc),
	}

	go func() {
		log.Printf("Starting server on port %s", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
