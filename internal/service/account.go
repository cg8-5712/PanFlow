package service

import (
	"context"
	"errors"
	"sync/atomic"

	"panflow/internal/model"
	"panflow/internal/repository"
)

var (
	ErrNoAvailableAccount = errors.New("no available account")
)

type AccountService struct {
	repo    *repository.AccountRepository
	client  *bdwpClient
	counter atomic.Uint64
}

func NewAccountService(repo *repository.AccountRepository, proxyURL string) *AccountService {
	return &AccountService{repo: repo, client: newBdwpClient(proxyURL)}
}

// PickForUser selects an available account for the given user.
// SVIP users only get their own account; others get the next account via round-robin.
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

	idx := s.counter.Add(1) % uint64(len(accounts))
	return &accounts[idx], nil
}

// RecordUsage increments account usage counters
func (s *AccountService) RecordUsage(ctx context.Context, id uint, size int64) error {
	return s.repo.IncrementUsage(id, size)
}

// CheckBanStatus calls the Baidu APL API to check if the account is banned/speed-limited
func (s *AccountService) CheckBanStatus(accountType, cookieOrToken, userAgent string, cid int64) (*BanStatus, error) {
	return s.client.CheckBanStatus(accountType, cookieOrToken, userAgent, cid)
}

// GetEnterpriseCID fetches the enterprise drive CID for the account identified by cookie
func (s *AccountService) GetEnterpriseCID(cookie, userAgent string) (int64, error) {
	return s.client.GetEnterpriseCID(cookie, userAgent)
}
