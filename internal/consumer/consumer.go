package consumer

import (
	"context"
	"encoding/json"

	"github.com/realdanielursul/order-service/internal/entity"
	"github.com/realdanielursul/order-service/internal/service"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func StartConsumer(service *service.Service, reader *kafka.Reader) {
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				logrus.Printf("kafka read error: %v\n", err)
				continue
			}

			var order entity.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				logrus.Printf("invalid message: %v\n", err)
				continue
			}

			if err := service.CreateOrder(context.Background(), &order); err != nil {
				logrus.Printf("failed to save order: %v", err)
			} else {
				logrus.Printf("order saved: %s", order.OrderUID)
			}
		}
	}()
}
