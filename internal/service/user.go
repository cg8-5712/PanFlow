package service

import (
	"context"
	"errors"

	"panflow/internal/model"
	"panflow/internal/repository"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDailyLimit    = errors.New("daily limit exceeded")
	ErrUserVipInsufficient = errors.New("vip balance insufficient")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetByID returns a user, checking cache first
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	key := UserCacheKey(id)

	var user model.User
	if CacheGet(ctx, key, &user) {
		return &user, nil
	}

	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	_ = CacheSet(ctx, key, u, ttlL1Short, ttlL2Medium)
	return u, nil
}

// CheckQuota verifies the user has remaining quota for a parse request
func (s *UserService) CheckQuota(ctx context.Context, userID uint) (*model.User, error) {
	u, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	switch u.UserType {
	case "admin":
		// unlimited
		return u, nil

	case "svip":
		if int64(u.DailyLimit) > 0 && u.DailyUsedCount >= int64(u.DailyLimit) {
			return nil, ErrUserDailyLimit
		}

	case "vip":
		if u.VipBalance <= 0 {
			return nil, ErrUserVipInsufficient
		}

	default: // guest
		if int64(u.DailyLimit) > 0 && u.DailyUsedCount >= int64(u.DailyLimit) {
			return nil, ErrUserDailyLimit
		}
	}

	return u, nil
}

// RecordUsage increments usage and deducts balance where applicable
func (s *UserService) RecordUsage(ctx context.Context, userID uint, userType string) error {
	if err := s.repo.IncrementDailyUsed(userID, 1); err != nil {
		return err
	}

	if userType == "vip" {
		if err := s.repo.DeductVipBalance(userID, 1); err != nil {
			return err
		}
	}

	// Invalidate cache
	CacheDelete(ctx, UserCacheKey(userID))
	return nil
}

// InvalidateCache removes a user from all cache layers
func (s *UserService) InvalidateCache(ctx context.Context, userID uint) {
	CacheDelete(ctx, UserCacheKey(userID))
}
