package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
			if errors.Is(err, context.Canceled) {
				return nil
			}

			fmt.Printf("consumer read failed: %v\n", err)

			time.Sleep(2 * time.Second)
			continue
		}

		var event AccountEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			return err
		}

		fmt.Printf("AUDIT: %s account=%s amount=%d\n", event.Type, event.AccountID, event.Amount)
	}
}
