package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/realdanielursul/order-service/config"
	"github.com/realdanielursul/order-service/internal/cache"
	"github.com/realdanielursul/order-service/internal/consumer"
	"github.com/realdanielursul/order-service/internal/handler"
	"github.com/realdanielursul/order-service/internal/repository"
	"github.com/realdanielursul/order-service/internal/service"
	"github.com/realdanielursul/order-service/pkg/httpserver"
	"github.com/realdanielursul/order-service/pkg/kafka"
	"github.com/realdanielursul/order-service/pkg/logger"
	"github.com/realdanielursul/order-service/pkg/postgres"
	"github.com/realdanielursul/order-service/pkg/redis"
	"github.com/sirupsen/logrus"
)

// ADD README + VIDEO

func main() {
	// Set up logger
	logger.SetLogrus()

	// Configure app
	cfg, err := config.NewConfig("./config/docker.yaml")
	if err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	// Connect to Redis Client
	client, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logrus.Fatalf("failed to connect to redis client: %s", err.Error())
	}

	// Connect to DB
	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("failed to connect to db: %s", err.Error())
	}

	// Initialize layers
	cache := cache.NewCache(client)
	repository := repository.NewRepository(db)
	service := service.NewService(cache, repository)

	// Preload cache
	if err := service.PreloadCache(context.Background()); err != nil {
		logrus.Fatalf("cache preload failed: %v", err)
	}

	// Connect to Kafka
	reader := kafka.NewKafkaReader(cfg.Kafka)

	// Start Kafka consumer
	consumer.StartConsumer(service, reader)

	// Run HTTP server
	handler := handler.NewHandler(service)
	srv := &httpserver.Server{}

	go func() {
		if err := srv.Run(cfg.HTTP.Port, handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error running http server: %s", err.Error())
		}
	}()

	logrus.Printf("App '%s %s' Started", cfg.App.Name, cfg.App.Version)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

	logrus.Printf("App '%s %s' Shutted Down", cfg.App.Name, cfg.App.Version)
}
