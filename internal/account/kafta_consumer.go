package account

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   "account-events",
			GroupID: "audit-service",
		}),
	}
}

func (consumer *KafkaConsumer) Consume(ctx context.Context) error {
	for {
		message, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event AccountEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			return err
		}

		fmt.Printf("AUDIT: %s account=%s amount=%d\n", event.Type, event.AccountID, event.Amount)
	}
}
