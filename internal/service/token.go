package service

import (
	"context"
	"errors"
	"time"

	"panflow/internal/model"
	"panflow/internal/repository"
)

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenDisabled   = errors.New("token is disabled")
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenQuotaCount = errors.New("token count quota exceeded")
	ErrTokenQuotaSize  = errors.New("token size quota exceeded")
	ErrTokenQuotaDay   = errors.New("token daily quota exceeded")
	ErrTokenIPLimit    = errors.New("token ip limit exceeded")
)

type TokenService struct {
	repo *repository.TokenRepository
}

func NewTokenService(repo *repository.TokenRepository) *TokenService {
	return &TokenService{repo: repo}
}

// GetByToken returns a token, checking L1→L2→DB
func (s *TokenService) GetByToken(ctx context.Context, tokenStr string) (*model.Token, error) {
	key := TokenCacheKey(tokenStr)

	var token model.Token
	if CacheGet(ctx, key, &token) {
		return &token, nil
	}

	t, err := s.repo.GetByToken(tokenStr)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	_ = CacheSet(ctx, key, t, ttlL1Short, ttlL2Medium)
	return t, nil
}

// Validate checks whether a token is valid and has quota for the given size
func (s *TokenService) Validate(ctx context.Context, tokenStr string, totalSize int64, clientIP string) (*model.Token, error) {
	t, err := s.GetByToken(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	if !t.Switch {
		return nil, ErrTokenDisabled
	}

	if t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Count quota (0 = unlimited)
	if t.Count > 0 && t.UsedCount >= t.Count {
		return nil, ErrTokenQuotaCount
	}

	// Size quota (0 = unlimited)
	if t.Size > 0 && t.UsedSize+totalSize > t.Size {
		return nil, ErrTokenQuotaSize
	}

	// Daily quota
	if t.TokenType == "daily" && t.Day > 0 {
		// daily reset is handled externally; UsedCount tracks today's usage
		if t.UsedCount >= t.Day {
			return nil, ErrTokenQuotaDay
		}
	}

	// IP limit
	if t.CanUseIPCount > 0 && clientIP != "" {
		found := false
		for _, ip := range t.IP {
			if ip == clientIP {
				found = true
				break
			}
		}
		if !found && int64(len(t.IP)) >= t.CanUseIPCount {
			return nil, ErrTokenIPLimit
		}
	}

	return t, nil
}

// RecordUsage increments usage counters and invalidates cache
func (s *TokenService) RecordUsage(ctx context.Context, id uint, size int64) error {
	if err := s.repo.IncrementUsage(id, size); err != nil {
		return err
	}
	// Invalidate cache so next read is fresh
	t, err := s.repo.GetByID(id)
	if err == nil {
		key := TokenCacheKey(t.Token)
		CacheDelete(ctx, key)
	}
	return nil
}

// InvalidateCache removes a token from all cache layers
func (s *TokenService) InvalidateCache(ctx context.Context, tokenStr string) {
	CacheDelete(ctx, TokenCacheKey(tokenStr))
}
