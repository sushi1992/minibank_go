package account

type AccountEvent struct {
	Type      string `json:"type"`
	AccountID string `json:"accountId"`
	Amount    int64  `json:"amount"`
}

const (
	EventAccountDeposited = "AccountDeposited"
	EventAccountWithdrawn = "AccountWithdrawn"
)

type EventPublisher interface {
	Publish(topic string, key string, event any) error
}
