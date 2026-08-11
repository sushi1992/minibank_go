package account

import "testing"

func TestTransactionProcessorProcessesDeposit(t *testing.T) {
	repo := NewMemoryRepository()

	account, err := NewAccount("Acc1", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	if err := repo.Save(account); err != nil {
		t.Fatalf("failed to save account: %v", err)
	}

	processor := NewTransactionProcessor(repo)

	go processor.Run()

	processor.Submit(Transaction{
		AccountID: "Acc1",
		Type:      TransactionDeposit,
		Amount:    500,
	})
}
