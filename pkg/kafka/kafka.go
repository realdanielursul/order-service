package kafka

import "github.com/segmentio/kafka-go"

func NewKafkaReader(broker, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
}
