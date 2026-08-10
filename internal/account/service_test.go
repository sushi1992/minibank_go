package account

import "testing"

func TestCreateAccountForService(t *testing.T) {
	id := "Apple"
	owner := "John Doe"
	service := NewService(NewMemoryRepository())
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

	savedAccount, err := service.repo.Get(id)
	if savedAccount != account {
		t.Fatalf("saved account should have been the same as the created account")
	}

	if err != nil {
		t.Fatalf("retrieving account caused an error: %v", err)
	}
}
