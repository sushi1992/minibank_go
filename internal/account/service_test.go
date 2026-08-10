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
