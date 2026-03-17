package cache

import (
	"context"
	"fmt"
	"time"

	"panflow/internal/config"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes the L2 distributed cache (Redis)
func InitRedis(cfg config.RedisConfig) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// RedisGet retrieves a value from L2 cache
func RedisGet(ctx context.Context, key string) (string, error) {
	if redisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	return redisClient.Get(ctx, key).Result()
}

// RedisSet stores a value in L2 cache with TTL
func RedisSet(ctx context.Context, key string, value any, ttl time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.Set(ctx, key, value, ttl).Err()
}

// RedisDelete removes a value from L2 cache
func RedisDelete(ctx context.Context, keys ...string) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.Del(ctx, keys...).Err()
}

// RedisExists checks if a key exists in L2 cache
func RedisExists(ctx context.Context, keys ...string) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}
	return redisClient.Exists(ctx, keys...).Result()
}

// RedisExpire sets a TTL on an existing key
func RedisExpire(ctx context.Context, key string, ttl time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.Expire(ctx, key, ttl).Err()
}

// RedisFlushDB clears all keys in the current database
func RedisFlushDB(ctx context.Context) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.FlushDB(ctx).Err()
}

// RedisClose closes the Redis connection
func RedisClose() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// RedisClient returns the underlying Redis client for advanced operations
func RedisClient() *redis.Client {
	return redisClient
}

// RedisHGet retrieves a field from a hash
func RedisHGet(ctx context.Context, key, field string) (string, error) {
	if redisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	return redisClient.HGet(ctx, key, field).Result()
}

// RedisHSet stores a field in a hash
func RedisHSet(ctx context.Context, key string, values ...any) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return redisClient.HSet(ctx, key, values...).Err()
}

// RedisHGetAll retrieves all fields from a hash
func RedisHGetAll(ctx context.Context, key string) (map[string]string, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}
	return redisClient.HGetAll(ctx, key).Result()
}

// RedisIncr increments a counter
func RedisIncr(ctx context.Context, key string) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}
	return redisClient.Incr(ctx, key).Result()
}

// RedisIncrBy increments a counter by a specific amount
func RedisIncrBy(ctx context.Context, key string, value int64) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}
	return redisClient.IncrBy(ctx, key, value).Result()
}
