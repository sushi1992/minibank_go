package account

import (
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type AccountEvent struct {
	Type      string     `json:"type"`
	AccountID string     `json:"accountId"`
	Amount    int64      `json:"amount"`
	EventID   gocql.UUID `json:"eventId"`
}

type OutboxEvent struct {
	CreatedAt time.Time
	EventID   gocql.UUID
	AccountID string
	Type      string
	Payload   string
}

const (
	EventAccountDeposited = "AccountDeposited"
	EventAccountWithdrawn = "AccountWithdrawn"
)

type EventPublisher interface {
	Publish(topic string, key string, event any) error
}
