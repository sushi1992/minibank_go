package account

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type OutboxProcessor struct {
	store     OutboxStore
	publisher EventPublisher
}

type OutboxStore interface {
	GetPendingOutboxEvents(limit int) ([]OutboxEvent, error)
	DeleteOutboxEvent(event OutboxEvent) error
}

func NewOutboxProcessor(store OutboxStore, publisher EventPublisher) *OutboxProcessor {
	return &OutboxProcessor{
		store:     store,
		publisher: publisher,
	}
}

func (processor *OutboxProcessor) Run(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := processor.processPendingEvents()
			if err != nil {
				fmt.Printf("failed processing outbox: %v\n", err)
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (processor *OutboxProcessor) processPendingEvents() error {
	events, err := processor.store.GetPendingOutboxEvents(100)
	if err != nil {
		return err
	}

	for _, event := range events {
		var accountEvent AccountEvent

		err := json.Unmarshal(
			[]byte(event.Payload),
			&accountEvent,
		)
		if err != nil {
			return err
		}

		err = processor.publisher.Publish(
			"account-events",
			event.AccountID,
			accountEvent,
		)

		if err != nil {
			return err
		}

		err = processor.store.DeleteOutboxEvent(event)
		if err != nil {
			return err
		}
	}

	return nil
}
