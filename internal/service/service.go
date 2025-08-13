package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/realdanielursul/order-service/internal/cache"
	"github.com/realdanielursul/order-service/internal/entity"
	"github.com/realdanielursul/order-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
		return nil, fmt.Errorf("get from cache: %w", err)
	}

	if data != nil {
		var order entity.Order
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, fmt.Errorf("unmarshal cached order: %w", err)
		}

		return &order, err
	}

	// get data from database
	order, err := s.repository.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("get from repository: %w", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	// set new data to cache
	data, err = json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("marshal new order: %w", err)
	}

	if err := s.cache.SetData(ctx, orderUID, data); err != nil {
		logrus.Warn("failed to cache order %q: %v", orderUID, err)
	}

	return order, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *entity.Order) error {
	// validate data

	// дисклеймер: не до конца понимаю, что означают все поля структуры
	// Order, следовательно не могу провалидировать входные данные,
	// поэтому будем считать, что здесь происходит **DATA VALIDATION**

	// save new data to database
	if err := s.repository.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("create order in repository: %w", err)
	}

	// set new data to cache
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal new order: %w", err)
	}

	if err := s.cache.SetData(ctx, order.OrderUID, data); err != nil {
		logrus.Warn("failed to cache order %q: %v", order.OrderUID, err)
	}

	return nil
}

func (s *Service) PreloadCache(ctx context.Context) error {
	orders, err := s.repository.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to preload cache: %w", err)
	}

	for _, order := range orders {
		data, err := json.Marshal(order)
		if err != nil {
			logrus.Warnf("failed to marshal order %q: %v", order.OrderUID, err)
			continue
		}

		if err := s.cache.SetData(ctx, order.OrderUID, data); err != nil {
			logrus.Warnf("failed to cache order %q: %v", order.OrderUID, err)
		}
	}

	logrus.Infof("Preloaded %d orders into cache", len(orders))
	return nil
}
