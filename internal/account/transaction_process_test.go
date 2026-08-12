package account

import (
	"context"
	"testing"
)

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
	result := make(chan TransactionResult)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.Run(ctx)

	err = processor.Submit(ctx, Transaction{
		AccountID: "Acc1",
		Type:      TransactionDeposit,
		Amount:    500,
		Result:    result,
	})

	if err != nil {
		t.Fatal("timeout occurred")
	}

	transactionResult := <-result

	if transactionResult.Err != nil {
		t.Fatalf("transaction failed: %v", transactionResult.Err)
	}

	if account.BalancePence != 500 {
		t.Fatalf("expected balance 500, got %d", account.BalancePence)
	}
}

func TestBufferedChannel(t *testing.T) {
	ch := make(chan int, 2)

	ch <- 10
	ch <- 20

	first := <-ch

	ch <- 30

	second := <-ch
	third := <-ch

	if first != 10 {
		t.Fatalf("expected 10, got %d", first)
	}

	if second != 20 {
		t.Fatalf("expected 20, got %d", second)
	}

	if third != 30 {
		t.Fatalf("expected 30, got %d", third)
	}
}

func TestBufferedChannelWithGoroutine(t *testing.T) {
	ch := make(chan int, 2)

	ch <- 10
	ch <- 20

	go func() {
		ch <- 30
	}()

	first := <-ch
	second := <-ch
	third := <-ch

	if first != 10 {
		t.Fatalf("expected 10, got %d", first)
	}

	if second != 20 {
		t.Fatalf("expected 20, got %d", second)
	}

	if third != 30 {
		t.Fatalf("expected 30, got %d", third)
	}
}

func TestTransactionProcessorWithdrawalFailure(t *testing.T) {
	// submit a withdrawal of 100 against an account with 0,
	// receive the result, and assert transactionResult.Err != nil.
	repo := NewMemoryRepository()

	accountId := "Acc1"
	account, err := NewAccount(accountId, "Owner", CurrencyGBP)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	if err := repo.Save(account); err != nil {
		t.Fatalf("failed to save account: %v", err)
	}

	processor := NewTransactionProcessor(repo)
	result := make(chan TransactionResult)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.Run(ctx)

	err = processor.Submit(ctx, Transaction{
		AccountID: accountId,
		Type:      TransactionWithdrawal,
		Amount:    100,
		Result:    result,
	})

	if err != nil {
		t.Fatal("timeout occurred")
	}

	transactionResult := <-result
	if transactionResult.Err == nil {
		t.Fatal("transaction passed when should have failed")
	}

	if account.BalancePence != 0 {
		t.Fatalf("expected balance to remain 0, got %d", account.BalancePence)
	}
}
