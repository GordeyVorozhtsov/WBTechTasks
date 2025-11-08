package main

import (
	"context"
	"eventbooker/internal/booker"
	"eventbooker/internal/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/dbpg"
)

func initDB() *dbpg.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)

	db, err := dbpg.New(
		dsn,
		[]string{},
		&dbpg.Options{
			MaxOpenConns:    25,
			MaxIdleConns:    25,
			ConnMaxLifetime: time.Hour,
		},
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}

func main() {
	// Подключаем базу данных
	db := initDB()
	bookerSvc := &booker.Service{DB: db}

	// Инициализируем таблицы
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := bookerSvc.InitDatabase(ctx); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database initialized")

	// Запускаем очистку просроченных броней
	go bookerSvc.StartCleanupScheduler(context.Background(), 30*time.Second)

	// сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.NewServer(bookerSvc),
	}

	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server exited gracefully")
}
