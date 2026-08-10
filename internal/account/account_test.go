package account

import "testing"

const startingBalance = 1000

func initialWithdrawTestSetup(t *testing.T) *Account {
	t.Helper()

	account, err := NewAccount("Apple", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatal("should be created successfully")
	}

	if account == nil {
		t.Fatal("account should have been created")
	}

	err = account.Deposit(startingBalance)
	if err != nil {
		t.Fatalf("expected deposit to be successful, instead got error %v", err)
	}

	if account.BalancePence != startingBalance {
		t.Fatalf("account balance should be %d, but is %d", startingBalance, account.BalancePence)
	}

	return account
}

func TestInvalidWithdrawal(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{"zero withdrawal", 0},
		{"negative withdrawal", -300},
		{"insufficient funds", 1001},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := initialWithdrawTestSetup(t)
			err := account.Withdraw(test.amount)
			if err == nil {
				t.Fatal("expected withdrawal to be unsuccessful")
			}

			if account.BalancePence != startingBalance {
				t.Errorf("account balance should be %d, but is %d", startingBalance, account.BalancePence)
			}
		})
	}
}

func TestWithdrawMoreThanBalanceFailsWithError(t *testing.T) {
	account := initialWithdrawTestSetup(t)
	err := account.Withdraw(1001)
	if err == nil {
		t.Fatal("expected withdrawal to be unsuccessful")
	}

	if account.BalancePence != startingBalance {
		t.Fatalf("account balance should be %d, but is %d", startingBalance, account.BalancePence)
	}
}

func TestValidCurrencyCreatesAccount(t *testing.T) {
	testID := "Apple"
	testOwner := "John Doe"
	account, err := NewAccount(testID, testOwner, CurrencyGBP)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if account == nil {
		t.Fatal("account should have been created")
	}

	if account.BalancePence != 0 {
		t.Errorf("expected balance 0, got %d", account.BalancePence)
	}

	if account.Currency != CurrencyGBP {
		t.Errorf("expected currency to be %s, instead got %s", CurrencyGBP, account.Currency)
	}

	if account.ID != "Apple" {
		t.Errorf("expected ID to be %v, instead got %v", testID, account.ID)
	}

	if account.Owner != "John Doe" {
		t.Errorf("expected owner to be %v, instead got %v", testOwner, account.Owner)
	}
}

func TestInvalidCurrencyRejectsAccountCreation(t *testing.T) {
	account, err := NewAccount("Apple", "John Doe", Currency("BANANA"))
	if err == nil {
		t.Fatal("expected error but received none")
	}

	if account != nil {
		t.Fatal("account should not have been created")
	}
}

func TestDepositOf500(t *testing.T) {
	account, err := NewAccount("Apple", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatal("should be created successfully")
	}

	if account == nil {
		t.Fatal("account should have been created")
	}

	err = account.Deposit(500)
	if err != nil {
		t.Fatalf("expected deposit to be successful, instead got error %v", err)
	}

	if account.BalancePence != 500 {
		t.Fatalf("account balance should be 500, but is %d", account.BalancePence)
	}
}

func TestDepositOf0ResultsInError(t *testing.T) {
	account, err := NewAccount("Apple", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatal("should be created successfully")
	}

	if account == nil {
		t.Fatal("account should have been created")
	}

	previousBalance := account.BalancePence
	err = account.Deposit(0)
	if err == nil {
		t.Error("expected deposit to be unsuccessful")
	}

	if account.BalancePence != previousBalance {
		t.Fatalf("account balance should be %d, but is %d", previousBalance, account.BalancePence)
	}
}

func TestDepositOfNegativeValueResultsInError(t *testing.T) {
	account, err := NewAccount("Apple", "John Doe", CurrencyGBP)
	if err != nil {
		t.Fatal("should be created successfully")
	}

	if account == nil {
		t.Fatal("account should have been created")
	}

	previousBalance := account.BalancePence
	err = account.Deposit(-10)
	if err == nil {
		t.Error("expected deposit to be unsuccessful")
	}

	if account.BalancePence != previousBalance {
		t.Fatalf("account balance should be %d, but is %d", previousBalance, account.BalancePence)
	}
}
