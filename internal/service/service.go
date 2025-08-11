package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/realdanielursul/order-service/internal/cache"
	"github.com/realdanielursul/order-service/internal/entity"
	"github.com/realdanielursul/order-service/internal/repository"
)

type Service struct {
	cache      *cache.Cache
	repository *repository.Repository
}

func NewService(c *cache.Cache, r *repository.Repository) *Service {
	return &Service{cache: c, repository: r}
}

func (s *Service) GetOrder(ctx context.Context, orderUID string) (*entity.Order, error) {
	cacheKey := fmt.Sprintf("order:%s", orderUID)
	data, err := s.cache.GetData(ctx, cacheKey)
	if err == nil && data != nil {
		var order entity.Order
		if err := json.Unmarshal(data, &order); err == nil {
			log.Println("GOT DATA FROM CACHE")
			return &order, err
		}
	}

	order, err := s.repository.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(order)
	if err == nil {
		s.cache.SetData(ctx, cacheKey, data)
	}

	return order, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *entity.Order) error {
	cacheKey := fmt.Sprintf("order:%s", order.OrderUID)
	log.Println(cacheKey)
	data, err := json.Marshal(order)
	if err == nil {
		log.Println(s.cache.SetData(ctx, cacheKey, data))
	}

	return s.repository.CreateOrder(ctx, order)
}
