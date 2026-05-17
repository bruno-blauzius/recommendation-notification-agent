package main

import (
	"fmt"
	"log"

	"github.com/project-go-sender-recommendation-agent/internal/adapters/primary/rabbitmq"
	"github.com/project-go-sender-recommendation-agent/internal/adapters/secondary/postgres"
	redisadapter "github.com/project-go-sender-recommendation-agent/internal/adapters/secondary/redis"
	"github.com/project-go-sender-recommendation-agent/internal/core/usecases"
	"github.com/project-go-sender-recommendation-agent/internal/infrastructure/config"
	"github.com/project-go-sender-recommendation-agent/internal/infrastructure/database"
)

func main() {
	fmt.Println("Hello World")

	cfg := config.Load()

	redisClient, err := database.NewRedisConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	recommendationRepo := postgres.NewRecommendationRepository(db)
	idempotencyRepo := redisadapter.NewIdempotencyRepository(redisClient)
	recommendationService := usecases.NewRecommendationService(recommendationRepo)

	handler := rabbitmq.NewSendRecommendationHandler(recommendationService, idempotencyRepo)

	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQDSN, cfg.RabbitMQQueue, handler)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	log.Println("Application started successfully")

	if err := consumer.Consume(); err != nil {
		log.Fatalf("consumer error: %v", err)
	}
}
