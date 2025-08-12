package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/realdanielursul/order-service/internal/cache"
	"github.com/realdanielursul/order-service/internal/entity"
	"github.com/realdanielursul/order-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	cache      *cache.Cache
	repository *repository.Repository
}

func NewService(c *cache.Cache, r *repository.Repository) *Service {
	return &Service{cache: c, repository: r}
}

func (s *Service) GetOrder(ctx context.Context, orderUID string) (*entity.Order, error) {
	// try to fetch data from cache
	data, err := s.cache.GetData(ctx, orderUID)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if data != nil {
		var order entity.Order
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, err
		}

		return &order, err
	}

	// get data from database
	order, err := s.repository.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	// set new data to cache
	data, err = json.Marshal(order)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetData(ctx, orderUID, data); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *entity.Order) error {
	// validate data
	// дисклеймер: не до конца понимаю, что значат поля структуры Order,
	// поэтому будем считать, что здесь происходит **DATA VALIDATION**

	// save new data to database
	if err := s.repository.CreateOrder(ctx, order); err != nil {
		return err
	}

	// set new data to cache
	data, err := json.Marshal(order)
	if err == nil {
		return s.cache.SetData(ctx, order.OrderUID, data)
	}

	return nil
}
