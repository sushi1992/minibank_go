package account

import "testing"

func TestSavingAccountIntoRepository(t *testing.T) {
	account, err := NewAccount("acc-001", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatalf("account should have been created but instead received error: %v", err)
	}

	repo := NewMemoryRepository()
	err = repo.Save(account)

	if err != nil {
		t.Fatalf("account should have been saved but instead received error %v", err)
	}

	retrievedAccount, err := repo.Get("acc-001")
	if err != nil {
		t.Fatalf("account should have been retrieved but instead received error %v", err)
	}

	if account != retrievedAccount {
		t.Fatalf("account retrieved should have been the same as account saved")
	}
}

func TestGettingNonExistentAccountFails(t *testing.T) {
	repo := NewMemoryRepository()
	retrievedAccount, err := repo.Get("BANANA")
	if err == nil {
		t.Fatal("account should not have been retrieved but instead was")
	}

	if retrievedAccount != nil {
		t.Fatalf("account retrieved should have been nil")
	}
}
