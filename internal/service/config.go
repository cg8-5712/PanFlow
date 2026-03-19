package service

import (
	"context"
	"strconv"

	"panflow/internal/model"
	"panflow/internal/repository"
)

type ConfigService struct {
	repo *repository.ConfigRepository
}

func NewConfigService(repo *repository.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

// Get returns a config value string, checking L1 cache first
func (s *ConfigService) Get(ctx context.Context, key string) (string, error) {
	cacheKey := ConfigCacheKey(key)

	var cfg model.Config
	if CacheGetL1Only(cacheKey, &cfg) {
		return cfg.Value, nil
	}

	c, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}

	_ = CacheSetL1Only(cacheKey, c, ttlL1Medium)
	return c.Value, nil
}

// GetInt returns a config value as int with a fallback default
func (s *ConfigService) GetInt(ctx context.Context, key string, def int) int {
	v, err := s.Get(ctx, key)
	if err != nil {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// GetBool returns a config value as bool with a fallback default
func (s *ConfigService) GetBool(ctx context.Context, key string, def bool) bool {
	v, err := s.Get(ctx, key)
	if err != nil {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
