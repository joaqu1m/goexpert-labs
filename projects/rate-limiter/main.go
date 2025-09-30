package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/config"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/middleware"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/ratelimiter"
	"github.com/joaqu1m/goexpert-labs/projects/rate-limiter/internal/storage"
)

type Response struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	redisStorage := storage.NewRedisStorage(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	ctx := context.Background()
	if err := redisStorage.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	rateLimiter := ratelimiter.New(redisStorage, cfg)

	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware.Handler(http.HandlerFunc(handler)))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := redisStorage.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exited")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Request successful - Rate limiter is working!",
		Time:    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}
