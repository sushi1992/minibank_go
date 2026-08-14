package account

import (
	"encoding/json"
	"errors"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type FakeRepository struct {
	savedAccount *Account
	savedOutbox  *OutboxEvent
	saveErr      error
}

func (repo *FakeRepository) Save(account *Account) error {
	repo.savedAccount = account
	return repo.saveErr
}

func (repo *FakeRepository) Get(id string) (*Account, error) {
	return nil, nil
}

func (repo *FakeRepository) SaveAccountAndOutboxEvent(
	account *Account,
	event OutboxEvent,
) error {
	repo.savedAccount = account
	repo.savedOutbox = &event

	return repo.saveErr
}

func TestDepositCreatesOutboxEvent(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	_, err := service.CreateAccount("Acc1", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	updatedAccount, err := service.Deposit("Acc1", 500)
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}

	if updatedAccount.BalancePence != 500 {
		t.Fatalf(
			"expected balance 500, got %d",
			updatedAccount.BalancePence,
		)
	}

	if len(repo.outbox) != 1 {
		t.Fatalf(
			"expected 1 outbox event, got %d",
			len(repo.outbox),
		)
	}

	event := repo.outbox[0]

	if event.AccountID != "Acc1" {
		t.Errorf(
			"expected account ID Acc1, got %s",
			event.AccountID,
		)
	}

	if event.Type != EventAccountDeposited {
		t.Errorf(
			"expected event type %s, got %s",
			EventAccountDeposited,
			event.Type,
		)
	}

	if event.EventID == (gocql.UUID{}) {
		t.Error("expected event ID to be populated")
	}

	if event.Payload == "" {
		t.Fatal("expected event payload to be populated")
	}

	var accountEvent AccountEvent

	err = json.Unmarshal([]byte(event.Payload), &accountEvent)
	if err != nil {
		t.Fatalf("failed to unmarshal outbox payload: %v", err)
	}

	if accountEvent.AccountID != "Acc1" {
		t.Errorf(
			"expected payload account ID Acc1, got %s",
			accountEvent.AccountID,
		)
	}

	if accountEvent.Amount != 500 {
		t.Errorf(
			"expected payload amount 500, got %d",
			accountEvent.Amount,
		)
	}

	if accountEvent.Type != EventAccountDeposited {
		t.Errorf(
			"expected payload event type %s, got %s",
			EventAccountDeposited,
			accountEvent.Type,
		)
	}
}

func TestCreateAccountForService(t *testing.T) {
	id := "Apple"
	owner := "John Doe"

	repo := &FakeRepository{}
	service := NewService(repo)
	account, err := service.CreateAccount(id, owner, CurrencyGBP)

	if err != nil {
		t.Fatalf("account creation should have been successful but received error: %v", err)
	}

	if account == nil {
		t.Fatal("expected account to be returned")
	}

	if account.ID != id {
		t.Errorf("expected '%s' but instead contains '%s'", id, account.ID)
	}

	if account.Owner != owner {
		t.Errorf("expected '%s' but instead contains '%s'", owner, account.Owner)
	}

	if account.Currency != CurrencyGBP {
		t.Errorf("expected '%s' but instead contains '%s'", CurrencyGBP, account.Currency)
	}

	if repo.savedAccount != account {
		t.Fatalf("saved account should have been the same as the created account")
	}
}

func TestCreateAccountFailsWhenSaveFailsForService(t *testing.T) {
	id := "Apple"
	owner := "John Doe"

	saveErr := errors.New("error saving account")
	repo := &FakeRepository{
		saveErr: saveErr,
	}

	service := NewService(repo)
	account, err := service.CreateAccount(id, owner, CurrencyGBP)

	if !errors.Is(err, saveErr) {
		t.Fatalf("expected original save error to be preserved, but instead got %v", err)
	}

	if account != nil {
		t.Fatalf("account should have been nil but instead is %v", account)
	}
}

func TestTransfer300Works(t *testing.T) {
	amountToTransfer := 300
	sourceId := "Acc1"
	owner := "John Doe"
	destinationId := "Acc2"

	repo := NewMemoryRepository()
	service := NewService(repo)
	sourceAccount, err := service.CreateAccount(sourceId, owner, CurrencyGBP)
	if err != nil {
		t.Fatalf("error occurred when creating account: %v", err)
	}
	err = sourceAccount.Deposit(1000)
	if err != nil {
		t.Fatalf("error funding source account: %v", err)
	}

	destinationAccount, err := service.CreateAccount(destinationId, owner, CurrencyGBP)
	if err != nil {
		t.Fatalf("error occurred when creating account: %v", err)
	}

	err = destinationAccount.Deposit(500)
	if err != nil {
		t.Fatalf("error funding destination account: %v", err)
	}

	err = service.Transfer(sourceId, destinationId, int64(amountToTransfer))
	if err != nil {
		t.Fatalf("error transferring amount: %v", err)
	}

	sourceBalance := sourceAccount.BalancePence
	if sourceBalance != 700 {
		t.Fatalf("source account has incorrect balance, should have been 700, instead it's %v", sourceBalance)
	}

	destinationBalance := destinationAccount.BalancePence
	if destinationBalance != 800 {
		t.Fatalf("source account has incorrect balance, should have been 800, instead it's %v", sourceBalance)
	}
}
