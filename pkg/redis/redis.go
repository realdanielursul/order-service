package redis

import (
	"context"

	"github.com/realdanielursul/order-service/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
