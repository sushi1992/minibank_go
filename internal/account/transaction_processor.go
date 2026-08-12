package account

import "context"

type TransactionType int

const (
	TransactionDeposit TransactionType = iota
	TransactionWithdrawal
)

type Transaction struct {
	AccountID string
	Type      TransactionType
	Amount    int64
	Result    chan TransactionResult
}

type TransactionResult struct {
	Err error
}

type TransactionProcessor struct {
	repo         Repository
	transactions chan Transaction
}

func NewTransactionProcessor(repo Repository) *TransactionProcessor {
	return &TransactionProcessor{
		repo:         repo,
		transactions: make(chan Transaction, 10),
	}
}

func (p *TransactionProcessor) Submit(ctx context.Context, transaction Transaction) error {
	select {
	case p.transactions <- transaction:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *TransactionProcessor) Run(ctx context.Context) {
	for {
		select {
		case transaction := <-p.transactions:
			account, err := p.repo.Get(transaction.AccountID)
			if err != nil {
				transaction.Result <- TransactionResult{Err: err}
				continue
			}

			switch transaction.Type {
			case TransactionDeposit:
				err = account.Deposit(transaction.Amount)
			case TransactionWithdrawal:
				err = account.Withdraw(transaction.Amount)
			}

			if err != nil {
				transaction.Result <- TransactionResult{Err: err}
				continue
			}

			err = p.repo.Save(account)
			if err != nil {
				transaction.Result <- TransactionResult{Err: err}
				continue
			}
			transaction.Result <- TransactionResult{Err: nil}

		case <-ctx.Done():
			return
		}
	}
}
