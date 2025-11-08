package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"imageprocessor/internal/server"
	"imageprocessor/internal/worker"

	kafka "github.com/wb-go/wbf/kafka"
)

var (
	PORT    = ":8080"
	BROKERS = []string{"localhost:9092"}
	TOPIC   = "image_tasks"
)

func main() {
	// создаём Kafka продюсер
	prod := kafka.NewProducer(BROKERS, TOPIC)

	// создаём Kafka консьюмер для воркера
	cons := kafka.NewConsumer(BROKERS, TOPIC, "imageprocessor-group")

	// создаём воркера
	w := worker.NewWorker(cons)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Start(ctx)

	srv := &http.Server{
		Addr:    PORT,
		Handler: server.NewServer(prod),
	}

	go func() {
		log.Printf("Server running on %s", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel() // останавливаем воркер

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// закрываем Kafka
	prod.Close()
	cons.Close()

	log.Println("Server exited")
}
