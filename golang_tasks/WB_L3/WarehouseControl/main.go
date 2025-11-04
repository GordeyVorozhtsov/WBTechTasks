package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	controller "warehousecontrol/internal/controller"
	"warehousecontrol/internal/server"
)

var PORT = ":8080"

func main() {
	// Инициализация сервиса (подключение к БД)
	svc := controller.NewService()
	// В данном примере Service.Close закрывает master/slave если нужно
	defer svc.Close()

	router := server.NewServer(svc)

	srv := &http.Server{
		Addr:         PORT,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Force shutdown: %v", err)
	}
	log.Println("Server stopped")
}
