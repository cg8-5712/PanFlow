package service

import (
	"context"
	"errors"
	"math/rand"

	"panflow/internal/model"
	"panflow/internal/repository"
)

var (
	ErrNoAvailableAccount = errors.New("no available account")
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// PickForUser selects an available account for the given user.
// SVIP users only get their own account; others get a random account from the pool.
func (s *AccountService) PickForUser(ctx context.Context, user *model.User) (*model.Account, error) {
	if user != nil && user.UserType == "svip" && user.BaiduAccountID != nil {
		acc, err := s.repo.GetByID(*user.BaiduAccountID)
		if err != nil {
			return nil, ErrNoAvailableAccount
		}
		if !acc.Switch {
			return nil, ErrNoAvailableAccount
		}
		return acc, nil
	}

	accounts, err := s.repo.ListEnabled("")
	if err != nil || len(accounts) == 0 {
		return nil, ErrNoAvailableAccount
	}

	return &accounts[rand.Intn(len(accounts))], nil
}

// RecordUsage increments account usage counters
func (s *AccountService) RecordUsage(ctx context.Context, id uint, size int64) error {
	return s.repo.IncrementUsage(id, size)
}
