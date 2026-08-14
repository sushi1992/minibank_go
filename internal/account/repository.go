package account

import "errors"

type Repository interface {
	Save(account *Account) error
	Get(id string) (*Account, error)
}

type AccountTransactionStore interface {
	SaveAccountAndOutboxEvent(
		account *Account,
		event OutboxEvent,
	) error
}

type AccountStore interface {
	Repository
	AccountTransactionStore
}

type MemoryRepository struct {
	accounts map[string]*Account
	outbox   []OutboxEvent
}

func (repo *MemoryRepository) SaveAccountAndOutboxEvent(
	account *Account,
	event OutboxEvent,
) error {
	repo.accounts[account.ID] = account
	repo.outbox = append(repo.outbox, event)

	return nil
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		accounts: make(map[string]*Account),
	}
}

func (repo *MemoryRepository) Save(account *Account) error {
	repo.accounts[account.ID] = account
	return nil
}

func (repo *MemoryRepository) Get(id string) (*Account, error) {
	account, ok := repo.accounts[id]
	if !ok {
		return nil, errors.New("account not found")
	}

	return account, nil
}
