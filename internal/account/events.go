package account

import gocql "github.com/apache/cassandra-gocql-driver/v2"

type AccountEvent struct {
	Type      string     `json:"type"`
	AccountID string     `json:"accountId"`
	Amount    int64      `json:"amount"`
	EventID   gocql.UUID `json:"eventId"`
}

type OutboxEvent struct {
	EventID   gocql.UUID
	AccountID string
	Type      string
	Payload   string
}

const (
	EventAccountDeposited = "AccountDeposited"
	EventAccountWithdrawn = "AccountWithdrawn"
)
