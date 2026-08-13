package account

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(broker string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr: kafka.TCP(broker),
		},
	}
}

func (publisher *KafkaPublisher) Publish(topic string, key string, event any) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: eventBytes,
	}
	return publisher.writer.WriteMessages(context.Background(), message)
}
