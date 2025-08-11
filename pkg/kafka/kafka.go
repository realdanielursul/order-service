package kafka

import (
	"github.com/realdanielursul/order-service/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaReader(cfg config.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Host + ":" + cfg.Port},
		Topic:   cfg.Topic,
	})
}
