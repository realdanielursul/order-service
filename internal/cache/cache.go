package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) SetData(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) GetData(ctx context.Context, key string) (string, error) {
	return c.Get(ctx, key).Result()
}

func (r *Cache) DeleteData(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
