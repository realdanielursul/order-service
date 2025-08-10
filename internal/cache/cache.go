package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = time.Hour

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) SetData(ctx context.Context, key string, value []byte) error {
	return c.Set(ctx, key, value, cacheTTL).Err()
}

func (c *Cache) GetData(ctx context.Context, key string) ([]byte, error) {
	res, err := c.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (r *Cache) DeleteData(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
