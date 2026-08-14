package account

import (
	"encoding/json"
	"fmt"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Service struct {
	store AccountStore
}

func NewService(store AccountStore) *Service {
	return &Service{
		store: store,
	}
}

func (service *Service) GetAccount(id string) (*Account, error) {
	account, err := service.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve account %w", err)
	}

	return account, nil
}

func (service *Service) Withdraw(id string, amount int64) (*Account, error) {
	account, err := service.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	err = account.Withdraw(amount)
	if err != nil {
		return nil, fmt.Errorf("withdrawing from account: %w", err)
	}

	err = service.store.Save(account)
	if err != nil {
		return nil, fmt.Errorf("saving account: %w", err)
	}

	return account, nil
}

func (service *Service) Deposit(id string, amount int64) (*Account, error) {
	account, err := service.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	err = account.Deposit(amount)
	if err != nil {
		return nil, fmt.Errorf("depositing into account: %w", err)
	}

	event := AccountEvent{
		Type:      EventAccountDeposited,
		AccountID: account.ID,
		Amount:    amount,
		EventID:   gocql.UUIDFromTime(time.Now()),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("json marshalling failed: %v", err)
	}

	outboxEvent := OutboxEvent{
		EventID:   event.EventID,
		AccountID: event.AccountID,
		Type:      event.Type,
		Payload:   string(payload),
	}

	err = service.store.SaveAccountAndOutboxEvent(account, outboxEvent)
	if err != nil {
		return nil, fmt.Errorf("publishing deposit event: %w", err)
	}
	return account, nil
}

func (service *Service) CreateAccount(id string, owner string, currency Currency) (*Account, error) {
	account, err := NewAccount(id, owner, currency)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	err = service.store.Save(account)
	if err != nil {
		return nil, fmt.Errorf("saving account: %w", err)
	}

	return account, nil
}

func (service *Service) Transfer(fromID string, toID string, amount int64) error {
	sourceAccount, err := service.store.Get(fromID)
	if err != nil {
		return fmt.Errorf("error retrieving source account: %w", err)
	}

	destinationAccount, err := service.store.Get(toID)
	if err != nil {
		return fmt.Errorf("error retrieving destination account: %w", err)
	}

	err = sourceAccount.Withdraw(amount)
	if err != nil {
		return fmt.Errorf("error occurred when attempting withdrawal: %w", err)
	}

	err = destinationAccount.Deposit(amount)
	if err != nil {
		return fmt.Errorf("error occurred when attempting deposit: %w", err)
	}

	err = service.store.Save(sourceAccount)
	if err != nil {
		return fmt.Errorf("error occurred when saving source account information: %w", err)
	}

	err = service.store.Save(destinationAccount)
	if err != nil {
		return fmt.Errorf("error occurred when saving destination account information: %w", err)
	}

	return nil
}
