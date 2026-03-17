package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"panflow/pkg/cache"
)

const (
	ttlL1Short  = 1 * time.Minute
	ttlL1Medium = 3 * time.Minute
	ttlL2Short  = 10 * time.Minute
	ttlL2Medium = 30 * time.Minute
)

// CacheGet retrieves a value from L1 then L2, unmarshalling into dest
func CacheGet(ctx context.Context, key string, dest any) bool {
	// L1
	if val, ok := cache.OtterGet(key); ok {
		if b, ok := val.([]byte); ok {
			if err := json.Unmarshal(b, dest); err == nil {
				return true
			}
		}
	}

	// L2
	raw, err := cache.RedisGet(ctx, key)
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		return false
	}

	// backfill L1
	cache.OtterSet(key, []byte(raw), ttlL1Medium)
	return true
}

// CacheSet stores a value in both L1 and L2
func CacheSet(ctx context.Context, key string, value any, l1TTL, l2TTL time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}

	cache.OtterSet(key, b, l1TTL)
	return cache.RedisSet(ctx, key, string(b), l2TTL)
}

// CacheDelete removes a key from both L1 and L2
func CacheDelete(ctx context.Context, key string) {
	cache.OtterDelete(key)
	_ = cache.RedisDelete(ctx, key)
}

// CacheSetL1Only stores a value only in L1 (for non-distributed data)
func CacheSetL1Only(key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}
	cache.OtterSet(key, b, ttl)
	return nil
}

// CacheGetL1Only retrieves a value only from L1
func CacheGetL1Only(key string, dest any) bool {
	val, ok := cache.OtterGet(key)
	if !ok {
		return false
	}
	b, ok := val.([]byte)
	if !ok {
		return false
	}
	return json.Unmarshal(b, dest) == nil
}

// ConfigCacheKey returns the cache key for a config entry
func ConfigCacheKey(key string) string {
	return "config:" + key
}

// TokenCacheKey returns the cache key for a token
func TokenCacheKey(token string) string {
	return "token:" + token
}

// UserCacheKey returns the cache key for a user
func UserCacheKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

// BlacklistCacheKey returns the cache key for a blacklist entry
func BlacklistCacheKey(typ, identifier string) string {
	return fmt.Sprintf("blacklist:%s:%s", typ, identifier)
}
