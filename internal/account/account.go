package account

import "errors"

type Account struct {
	ID           string
	Owner        string
	BalancePence int64
	Currency     Currency
}

type Currency string

const (
	CurrencyGBP Currency = "GBP"
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
)

func isValidCurrency(currency Currency) bool {
	switch currency {
	case CurrencyGBP, CurrencyEUR, CurrencyUSD:
		return true
	default:
		return false
	}
}

func NewAccount(id, owner string, currency Currency) (*Account, error) {
	if !isValidCurrency(currency) {
		return nil, errors.New("invalid currency")
	}

	account := Account{
		BalancePence: 0,
		ID:           id,
		Owner:        owner,
		Currency:     currency,
	}
	return &account, nil
}

func (account *Account) Deposit(amountToDeposit int64) error {
	if amountToDeposit <= 0 {
		return errors.New("amount to deposit should be greater than 0")
	}

	account.BalancePence += amountToDeposit
	return nil
}

func (account *Account) Withdraw(amountToWithdraw int64) error {
	if amountToWithdraw <= 0 {
		return errors.New("withdrawal amount must be greater than 0")
	}

	if amountToWithdraw > account.BalancePence {
		return errors.New("insufficient funds")
	}

	account.BalancePence -= amountToWithdraw
	return nil
}
