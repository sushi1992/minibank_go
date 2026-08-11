package account

import (
	"errors"
	"testing"
)

type FakeRepository struct {
	savedAccount *Account
	saveErr      error
}

func (repo *FakeRepository) Save(account *Account) error {
	repo.savedAccount = account
	return repo.saveErr
}

func (repo *FakeRepository) Get(id string) (*Account, error) {
	return nil, nil
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
