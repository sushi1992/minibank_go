package account

import (
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (service *Service) CreateAccount(id string, owner string, currency Currency) (*Account, error) {
	account, err := NewAccount(id, owner, currency)
	if err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	err = service.repo.Save(account)
	if err != nil {
		return nil, fmt.Errorf("saving account: %w", err)
	}

	return account, nil
}
